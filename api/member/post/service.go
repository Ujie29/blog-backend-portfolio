package post

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"

	"blog-backend/common/entity"
	"blog-backend/common/middleware"
	"blog-backend/common/model"

	"blog-backend/common/utils"

	"github.com/uptrace/bun"
)

type PostService interface {
	GetPostList(req GetPostListDto) (model.PaginatedResponse[PostListDto], error)
	GetPostBySlug(slug string) (PostDto, error)
	GetPostsByCategory(slug string, req GetPostListDto) (model.PaginatedResponse[PostListDto], error)
	GetAboutMe() (AboutMeDto, error)
	GetRandomPostsByCategory(dto GetRandomPostsByCategoryDto) ([]PostListDto, error)
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

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 15
	}

	// 查詢總筆數
	total, err := s.db.NewSelect().
		Model((*entity.Post)(nil)).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Count(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// 查詢分頁資料
	var posts []entity.Post
	err = s.db.NewSelect().
		Model(&posts).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Order("created_at DESC").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Scan(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// 組裝結果
	var result []PostListDto
	for _, post := range posts {
		result = append(result, PostListDto{
			Slug:          post.Slug,
			Title:         post.Title,
			Summary:       utils.ExtractSummaryFromEditorJS(post.Content, 200),
			CoverImageUrl: post.CoverImageUrl,
			CreatedAt:     post.CreatedAt,
		})
	}

	// 回傳文章清單、總筆數
	return model.PaginatedResponse[PostListDto]{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalCount: total,
		Data:       result,
	}, nil
}

func (s *postServiceImpl) GetPostBySlug(slug string) (PostDto, error) {
	var post entity.Post
	err := s.db.NewSelect().
		Model(&post).
		Where("post.slug = ?", slug).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Limit(1).
		Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return PostDto{}, middleware.ErrNotFound
	} else if err != nil {
		return PostDto{}, middleware.ErrDB
	}

	dto := PostDto{
		Title:         post.Title,
		Summary:       utils.ExtractSummaryFromEditorJS(post.Content, 200),
		Content:       post.Content,
		CategoryID:    post.CategoryID,
		CoverImageUrl: post.CoverImageUrl,
		CreatedAt:     post.CreatedAt,
	}
	return dto, nil
}

func (s *postServiceImpl) GetPostsByCategory(slug string, req GetPostListDto) (model.PaginatedResponse[PostListDto], error) {
	ctx := context.Background()

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 15
	}

	// 找出該分類是主分類還是子分類
	var category entity.Category
	err := s.db.NewSelect().
		Model(&category).
		Where("slug = ?", slug).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	var categoryIDs []uint

	if category.HasChildren {
		// 有子分類：查詢所有子分類 ID
		var subCategories []entity.Category
		err := s.db.NewSelect().
			Model(&subCategories).
			Where("parent = ?", category.ID).
			Scan(ctx)
		if err != nil {
			return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
		}

		for _, sub := range subCategories {
			categoryIDs = append(categoryIDs, sub.ID)
		}
	} else {
		// 沒有子分類：表示主分類可擁有自己的文章
		categoryIDs = append(categoryIDs, category.ID)
	}

	// ✅ 查詢總筆數
	total, err := s.db.NewSelect().
		Model((*entity.Post)(nil)).
		Where("category_id IN (?)", bun.In(categoryIDs)).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Count(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// ✅ 分頁查詢資料
	var posts []entity.Post
	err = s.db.NewSelect().
		Model(&posts).
		Where("category_id IN (?)", bun.In(categoryIDs)).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Order("created_at DESC").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Scan(ctx)
	if err != nil {
		return model.PaginatedResponse[PostListDto]{}, middleware.ErrDB
	}

	// 組裝 DTO
	var result []PostListDto
	for _, post := range posts {
		result = append(result, PostListDto{
			Slug:          post.Slug,
			Title:         post.Title,
			Summary:       utils.ExtractSummaryFromEditorJS(post.Content, 200),
			CoverImageUrl: post.CoverImageUrl,
			CreatedAt:     post.CreatedAt,
		})
	}

	return model.PaginatedResponse[PostListDto]{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalCount: total,
		Data:       result,
	}, nil
}

func (s *postServiceImpl) GetAboutMe() (AboutMeDto, error) {
	var about entity.AboutMe
	err := s.db.NewSelect().
		Model(&about).
		Order("updated_at DESC").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return AboutMeDto{}, middleware.ErrDB
	}

	return AboutMeDto{
		ID:        about.ID,
		Content:   about.HtmlContent,
		UpdatedAt: about.UpdatedAt,
	}, nil
}

func (s *postServiceImpl) GetRandomPostsByCategory(dto GetRandomPostsByCategoryDto) ([]PostListDto, error) {
	ctx := context.Background()

	var category entity.Category
	err := s.db.NewSelect().
		Model(&category).
		Where("id = ?", dto.CategoryID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, middleware.ErrDB
	}

	var categoryIDs []uint

	if category.HasChildren {
		// 有子分類 → 查找子分類的文章
		var subCategories []entity.Category
		err := s.db.NewSelect().
			Model(&subCategories).
			Where("parent = ?", category.ID).
			Scan(ctx)
		if err != nil {
			return nil, middleware.ErrDB
		}

		for _, sub := range subCategories {
			categoryIDs = append(categoryIDs, sub.ID)
		}
	} else {
		// 沒子分類 → 查自己的文章
		categoryIDs = append(categoryIDs, category.ID)
	}

	// ✅ 查詢文章，排除 slug
	var posts []entity.Post
	query := s.db.NewSelect().
		Model(&posts).
		Where("is_published = TRUE").
		Where("is_deleted = FALSE").
		Where("category_id IN (?)", bun.In(categoryIDs))

	if dto.Slug != "" {
		query = query.Where("slug != ?", dto.Slug)
	}

	err = query.Scan(ctx)
	if err != nil {
		return nil, middleware.ErrDB
	}

	if len(posts) == 0 {
		return []PostListDto{}, nil
	}

	// ✅ 隨機打亂
	for i := len(posts) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		posts[i], posts[j] = posts[j], posts[i]
	}

	// ✅ 取前六筆
	selected := posts
	if len(posts) > 6 {
		selected = posts[:6]
	}

	var result []PostListDto
	for _, post := range selected {
		result = append(result, PostListDto{
			Slug:          post.Slug,
			Title:         post.Title,
			Summary:       utils.ExtractSummaryFromEditorJS(post.Content, 200),
			CoverImageUrl: post.CoverImageUrl,
			CreatedAt:     post.CreatedAt,
		})
	}

	return result, nil
}
