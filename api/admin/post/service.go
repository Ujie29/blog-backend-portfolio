package post

import (
	"blog-backend/common/middleware"
	"blog-backend/common/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"blog-backend/common/entity"
	"blog-backend/common/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PostService interface {
	GetPostList(req GetPostListDto) (model.PaginatedResponse[PostListDto], error)
	GetPostByID(id string) (PostDto, error)
	CreatePost(req CreatePostDto) (PostDto, error)
	UpdatePost(id string, req UpdatePostDto) (PostDto, error)
	DeletePost(id string) error
	GetPostsByCategory(categoryID string) ([]PostListDto, error)
	GeneratePresignedUploadURL(filename string) (UploadUrlDto, error)
	GetAboutMe() (AboutMeDto, error)
	UpdateAboutMe(req UpdateAboutMeDto) (AboutMeDto, error)
}

type postServiceImpl struct {
	db *bun.DB
}

func NewPostService(db *bun.DB) PostService {
	return &postServiceImpl{
		db: db,
	}
}

func (s *postServiceImpl) GetPostList(req GetPostListDto) (model.PaginatedResponse[PostListDto], error) {
	ctx := context.Background()

	// 查詢總筆數
	countQuery := s.db.NewSelect().Model(&entity.Post{}).Where("is_deleted = false")
	if req.Search != "" {
		countQuery = countQuery.Where("title ILIKE ?", "%"+req.Search+"%")
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// 查詢分頁資料
	var posts []entity.Post
	query := s.db.NewSelect().Model(&posts).Where("is_deleted = false")

	if req.Search != "" {
		query = query.Where("title ILIKE ?", "%"+req.Search+"%")
	}

	err = query.
		Order("created_at DESC").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Scan(ctx)

	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// 組裝回傳 DTO
	var result []PostListDto
	for i, post := range posts {
		result = append(result, PostListDto{
			SortID:    (req.Page-1)*req.Limit + i + 1,
			Id:        post.ID,
			Title:     post.Title,
			Summary:   utils.ExtractSummaryFromEditorJS(post.Content, 200),
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	// 回傳統一分頁格式
	return model.PaginatedResponse[PostListDto]{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalCount: total,
		Data:       result,
	}, nil
}

func (s *postServiceImpl) GetPostByID(id string) (PostDto, error) {
	var post entity.Post
	err := s.db.NewSelect().
		Model(&post).
		Where("post.id = ?", id).
		Where("is_deleted = false").
		Limit(1).
		Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return PostDto{}, middleware.ErrNotFound
	} else if err != nil {
		return PostDto{}, middleware.ErrDB
	}
	dto := PostDto{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Summary:       utils.ExtractSummaryFromEditorJS(post.Content, 200),
		CoverImageUrl: post.CoverImageUrl,
		IsPublished:   post.IsPublished,
		CategoryID:    post.CategoryID,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		Slug:          post.Slug,
	}
	return dto, nil
}

func (s *postServiceImpl) CreatePost(req CreatePostDto) (PostDto, error) {
	now := time.Now()
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return PostDto{}, middleware.ErrTransaction
	}
	defer tx.Rollback()

	if strings.TrimSpace(req.Content) == "" || req.Content == "null" {
		return PostDto{}, middleware.ErrContentEmpty
	}

	// 建立 post
	post := entity.Post{
		Title:         req.Title,
		CategoryID:    req.CategoryID,
		IsPublished:   req.IsPublished,
		CoverImageUrl: req.CoverImageUrl,
		Content:       req.Content,
		CreatedAt:     now,
		UpdatedAt:     now,
		Slug:          req.Slug,
		NeedsRefresh:  true,
	}
	_, err = tx.NewInsert().Model(&post).Returning("id").Exec(ctx)
	if err != nil {
		return PostDto{}, middleware.ErrDB
	}

	// 處理圖片紀錄（從內容中擷取 image url）
	imageUrls := extractImageUrls(req.Content)
	var images []entity.Image
	for _, url := range imageUrls {
		images = append(images, entity.Image{
			URL:       url,
			PostID:    post.ID,
			Type:      "inline",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if req.CoverImageUrl != "" {
		images = append(images, entity.Image{
			URL:       req.CoverImageUrl,
			PostID:    post.ID,
			Type:      "cover",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if len(images) > 0 {
		_, err := tx.NewInsert().Model(&images).Exec(ctx)
		if err != nil {
			return PostDto{}, middleware.ErrDB
		}
	}

	if err := tx.Commit(); err != nil {
		return PostDto{}, middleware.ErrTransaction
	}

	// ✅ 清除快取 + 重新部署（不影響主流程）
	if req.IsPublished {
		go func() {
			if err := utils.PurgeWorkerCacheAndDeployVercel(); err != nil {
				fmt.Printf("⚠️ 部署失敗（CreatePost）：%v\n", err)
			}
		}()
	}

	return s.GetPostByID(fmt.Sprint(post.ID))
}

func (s *postServiceImpl) UpdatePost(id string, req UpdatePostDto) (PostDto, error) {
	ctx := context.Background()

	// 檢查空內容
	if strings.TrimSpace(req.Content) == "" || req.Content == "null" {
		return PostDto{}, middleware.ErrContentEmpty
	}

	// 取得原本文章與 content block
	var post entity.Post
	err := s.db.NewSelect().
		Model(&post).
		Where("post.id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return PostDto{}, middleware.ErrDB
	}

	// 比對新舊圖片 URL
	oldUrls := extractImageUrls(post.Content)
	newUrls := extractImageUrls(req.Content)
	added, removed := diffImageUrls(oldUrls, newUrls)
	oldCover := post.CoverImageUrl
	newCover := req.CoverImageUrl
	coverChanged := oldCover != "" && oldCover != newCover

	// 開始 Transaction 更新文章與內容
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return PostDto{}, middleware.ErrTransaction
	}
	defer tx.Rollback()

	// 更新主文章
	_, err = tx.NewUpdate().
		Model(&entity.Post{
			ID:            post.ID,
			Title:         req.Title,
			CategoryID:    req.CategoryID,
			IsPublished:   req.IsPublished,
			CoverImageUrl: req.CoverImageUrl,
			Content:       req.Content,
			CreatedAt:     post.CreatedAt,
			UpdatedAt:     time.Now(),
			Slug:          req.Slug,
			NeedsRefresh:  true,
		}).
		WherePK().
		Exec(ctx)
	if err != nil {
		return PostDto{}, middleware.ErrDB
	}

	// 新增新增的圖片
	for _, url := range added {
		_, _ = tx.NewInsert().
			Model(&entity.Image{
				URL:    url,
				PostID: post.ID,
				Type:   "inline",
				Status: "active",
			}).
			Ignore().
			Exec(ctx)
	}

	// 標記被移除的圖片
	if len(removed) > 0 {
		_, _ = tx.NewUpdate().
			Model((*entity.Image)(nil)).
			Set("status = 'pending_delete', updated_at = NOW()").
			Where("post_id = ?", post.ID).
			Where("url IN (?)", bun.In(removed)).
			Exec(ctx)
	}

	if coverChanged {
		_, _ = tx.NewUpdate().
			Model((*entity.Image)(nil)).
			Set("status = 'pending_delete', updated_at = NOW()").
			Where("post_id = ?", post.ID).
			Where("url = ? AND type = 'cover'", oldCover).
			Exec(ctx)

		if newCover != "" {
			_, _ = tx.NewInsert().
				Model(&entity.Image{
					URL:       newCover,
					PostID:    post.ID,
					Type:      "cover",
					Status:    "active",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}).
				Ignore().
				Exec(ctx)
		}
	}

	// 成功提交
	if err := tx.Commit(); err != nil {
		return PostDto{}, middleware.ErrTransaction
	}

	// ✅ 清除快取 + 重新部署（不影響主流程）
	if req.IsPublished {
		go func() {
			if err := utils.PurgeWorkerCacheAndDeployVercel(); err != nil {
				fmt.Printf("⚠️ 部署失敗（CreatePost）：%v\n", err)
			}
		}()
	}

	return s.GetPostByID(fmt.Sprint(post.ID))
}

func (s *postServiceImpl) DeletePost(id string) error {
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return middleware.ErrTransaction
	}
	defer tx.Rollback()

	// 將 images.status 改為 pending_delete
	_, err = tx.NewUpdate().
		Model((*entity.Image)(nil)).
		Set("status = 'pending_delete', updated_at = NOW()").
		Where("post_id = ?", id).
		Exec(ctx)
	if err != nil {
		return middleware.ErrDB
	}

	// 軟刪除：將 is_deleted 設為 true
	_, err = tx.NewUpdate().
		Model((*entity.Post)(nil)).
		Set("is_deleted = true, updated_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return middleware.ErrDB
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// ✅ 清除快取 + 重新部署（不影響主流程）
	go func() {
		if err := utils.PurgeWorkerCacheAndDeployVercel(); err != nil {
			fmt.Printf("⚠️ 部署失敗（DeletePost）：%v\n", err)
		}
	}()

	return nil
}

func (s *postServiceImpl) GetPostsByCategory(categoryID string) ([]PostListDto, error) {
	var posts []entity.Post
	err := s.db.NewSelect().
		Model(&posts).
		Where("category_id = ?", categoryID).
		Where("is_deleted = false").
		Scan(context.Background())
	if err != nil {
		return nil, middleware.ErrDB
	}

	var result []PostListDto = make([]PostListDto, 0)
	for _, post := range posts {
		dto := PostListDto{
			Title:   post.Title,
			Summary: utils.ExtractSummaryFromEditorJS(post.Content, 200),
		}
		result = append(result, dto)
	}

	return result, nil
}

func (s *postServiceImpl) GeneratePresignedUploadURL(filename string) (UploadUrlDto, error) {
	// ⚙️ 設定 AWS/R2 資訊
	bucket := "images" // R2 bucket 名稱
	region := "auto"   // R2 可用 "auto"

	endpoint := os.Getenv("R2_ENDPOINT")             // ✅ R2 endpoint
	accessKey := os.Getenv("R2_ACCESS_KEY")          // ✅ R2 Access Key
	secretKey := os.Getenv("R2_SECRET_KEY")          // ✅ R2 Secret Key
	publicBaseURL := os.Getenv("R2_PUBLIC_BASE_URL") // ✅ 對外公開的圖片 Base URL

	if endpoint == "" || accessKey == "" || secretKey == "" || publicBaseURL == "" {
		return UploadUrlDto{}, middleware.ErrExternalService
	}

	ext := filepath.Ext(filename) // 取得原始副檔名（如 .png）
	if ext == "" {
		ext = ".jpg" // 預設副檔名（防止空白）
	}
	newFilename := uuid.New().String() + ext

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		return UploadUrlDto{}, middleware.ErrExternalService
	}

	s3Client := s3.New(sess)

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(newFilename),
	})

	expiresIn := int64(300) // 預簽名 URL 的有效時間（秒）
	uploadUrl, err := req.Presign(time.Duration(expiresIn) * time.Second)
	if err != nil {
		return UploadUrlDto{}, middleware.ErrExternalService
	}

	return UploadUrlDto{
		UploadUrl: uploadUrl,
		ImageUrl:  fmt.Sprintf("%s/%s", publicBaseURL, newFilename),
		ExpiresIn: expiresIn,
		Filename:  newFilename,
	}, nil
}

func (s *postServiceImpl) GetAboutMe() (AboutMeDto, error) {
	var about entity.AboutMe
	err := s.db.NewSelect().
		Model(&about).
		Order("updated_at DESC").
		Limit(1).
		Scan(context.Background())

	// ✅ 如果找不到資料，就回傳一筆空的預設資料（不報錯）
	if errors.Is(err, sql.ErrNoRows) {
		return AboutMeDto{
			ID:        0,
			Content:   `{"time":0,"blocks":[],"version":"2.28.2"}`, // 空 Editor.js 資料
			UpdatedAt: time.Now(),
		}, nil
	} else if err != nil {
		// ✅ 其他錯誤才真的回傳錯
		return AboutMeDto{}, middleware.ErrDB
	}

	return AboutMeDto{
		ID:        about.ID,
		Content:   about.HtmlContent,
		UpdatedAt: about.UpdatedAt,
	}, nil
}

func (s *postServiceImpl) UpdateAboutMe(req UpdateAboutMeDto) (AboutMeDto, error) {
	ctx := context.Background()
	now := time.Now()

	var existing entity.AboutMe
	err := s.db.NewSelect().
		Model(&existing).
		Limit(1).
		Scan(ctx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AboutMeDto{}, middleware.ErrDB
	}

	oldContent := existing.HtmlContent
	oldUrls := extractImageUrls(oldContent)
	newUrls := extractImageUrls(req.Content)
	added, removed := diffImageUrls(oldUrls, newUrls)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return AboutMeDto{}, middleware.ErrDB
	}
	defer tx.Rollback()

	// 有舊資料 → 更新
	if existing.ID != 0 {
		existing.HtmlContent = req.Content
		existing.UpdatedAt = now

		_, err := tx.NewUpdate().
			Model(&existing).
			WherePK().
			Exec(ctx)
		if err != nil {
			return AboutMeDto{}, middleware.ErrDB
		}
	} else {
		// 沒資料 → 新增
		existing = entity.AboutMe{
			HtmlContent: req.Content,
			UpdatedAt:   now,
		}
		_, err := tx.NewInsert().
			Model(&existing).
			Exec(ctx)
		if err != nil {
			return AboutMeDto{}, middleware.ErrDB
		}
	}

	// ✅ 移除的圖片 → 標記 pending_delete
	if len(removed) > 0 {
		_, _ = tx.NewUpdate().
			Model((*entity.Image)(nil)).
			Set("status = 'pending_delete', updated_at = NOW()").
			Where("type = 'about'").
			Where("url IN (?)", bun.In(removed)).
			Exec(ctx)
	}

	// ✅ 新增的圖片
	for _, url := range added {
		_, _ = tx.NewInsert().
			Model(&entity.Image{
				URL:       url,
				Type:      "about",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: now,
			}).
			Ignore().
			Exec(ctx)
	}

	if err := tx.Commit(); err != nil {
		return AboutMeDto{}, middleware.Newf(middleware.ErrDB.Code, "提交交易失敗：%v", err)
	}

	// ✅ 清除快取 + 重新部署（不影響主流程）
	go func() {
		if err := utils.PurgeWorkerCacheAndDeployVercel(); err != nil {
			fmt.Printf("⚠️ 部署失敗（CreatePost）：%v\n", err)
		}
	}()

	return AboutMeDto{
		ID:        existing.ID,
		Content:   existing.HtmlContent,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

// ✅ 從文章內容中抓出所有圖片 URL（Editor.js JSON 解析）
func extractImageUrls(content string) []string {
	urls := []string{}
	if strings.TrimSpace(content) == "" {
		return urls
	}
	// 簡單抓取 `"file": { "url": "..." }`
	pattern := `"file":\s*{\s*"url":\s*"(.*?)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}
	return urls
}

// ✅ 新舊圖片 URL 差異比對：哪些是新增？哪些是被移除？
func diffImageUrls(oldUrls, newUrls []string) (added, removed []string) {
	oldSet := make(map[string]bool)
	newSet := make(map[string]bool)

	for _, u := range oldUrls {
		oldSet[u] = true
	}
	for _, u := range newUrls {
		newSet[u] = true
	}

	for _, u := range newUrls {
		if !oldSet[u] {
			added = append(added, u)
		}
	}
	for _, u := range oldUrls {
		if !newSet[u] {
			removed = append(removed, u)
		}
	}
	return
}
