package batch

import (
	"blog-backend/common/entity"
	"blog-backend/common/middleware"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/uptrace/bun"
)

type BatchService interface {
	CleanPendingImages() (int, error)
}

type batchServiceImpl struct {
	db *bun.DB
}

func NewBatchService(db *bun.DB) BatchService {
	return &batchServiceImpl{db: db}
}

func (s *batchServiceImpl) CleanPendingImages() (int, error) {
	ctx := context.Background()

	fmt.Println("🚀 開始執行圖片清理任務...")

	// 查出所有 status = pending_delete 且尚未刪除的圖片
	var images []entity.Image
	err := s.db.NewSelect().
		Model(&images).
		Where("status = 'pending_delete'").
		Where("is_deleted = FALSE").
		Scan(ctx)
	if err != nil {
		return 0, middleware.WrapDBErr("查詢待刪圖片失敗", err)
	}
	fmt.Printf("🔍 找到 %d 張待刪圖片\n", len(images))

	if len(images) == 0 {
		fmt.Println("✅ 無待刪圖片，結束任務")
		return 0, nil
	}

	// 初始化 S3 / R2 client
	s3Client, err := newR2Client()
	if err != nil {
		return 0, middleware.New("R2_INIT_FAILED", fmt.Sprintf("初始化 R2 client 失敗：%v", err))
	}

	bucket := "images"
	deletedCount := 0

	for i, img := range images {
		fmt.Printf("👉 [%d/%d] 處理圖片 ID=%s\n", i+1, len(images), img.ID)

		key := extractR2ObjectKey(img.URL)
		if key == "" {
			fmt.Printf("⚠️ 無法解析 R2 Key，跳過：%s\n", img.URL)
			continue
		}

		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			fmt.Printf("❌ 刪除失敗（R2）: ID=%s，Key=%s，錯誤: %v\n", img.ID, key, err)
			continue
		}

		// 軟刪除（is_deleted = true, deleted_at = NOW()）
		_, err = s.db.NewUpdate().
			Model(&entity.Image{}).
			Set("is_deleted = TRUE, deleted_at = NOW()").
			Where("id = ?", img.ID).
			Exec(ctx)
		if err != nil {
			fmt.Printf("⚠️ 資料庫刪除標記失敗: ID=%s，錯誤: %v\n", img.ID, err)
			continue
		}

		deletedCount++
		fmt.Printf("✅ 刪除成功：ID=%s，Key=%s\n", img.ID, key)
	}

	fmt.Printf("🎉 圖片清理任務完成，共成功刪除 %d 張圖片\n", deletedCount)

	return deletedCount, nil
}

// 建立 R2 client
func newR2Client() (*s3.S3, error) {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	region := "auto"

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("R2 環境變數未設定完整")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			accessKey, secretKey, "",
		),
	})
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

// 從公開網址中取出 R2 Object Key（給 R2 刪除用）
func extractR2ObjectKey(imageURL string) string {
	publicBase := os.Getenv("R2_PUBLIC_BASE_URL")
	if publicBase == "" {
		return ""
	}

	u, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}

	if !strings.HasPrefix(imageURL, publicBase) {
		return ""
	}

	return strings.TrimPrefix(u.Path, "/")
}
