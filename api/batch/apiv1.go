package batch

import (
	"blog-backend/common/middleware"

	"github.com/gin-gonic/gin"
)

// BatchAPI 是 Batch API 控制器
type BatchAPI struct {
	service BatchService
}

// NewBatchAPI 建立 BatchAPI 實例
func NewBatchAPI(service BatchService) *BatchAPI {
	return &BatchAPI{
		service: service,
	}
}

// RegisterRoutes 註冊路由
func (api *BatchAPI) RegisterRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api/batch")
	{
		apiGroup.POST("/clean-images", api.CleanPendingImages)
	}
}

// CleanPendingImages 清除 pending_delete 狀態的圖片（R2 + 資料庫）
func (api *BatchAPI) CleanPendingImages(c *gin.Context) {
	count, err := api.service.CleanPendingImages()
	if err != nil {
		c.Error(middleware.WrapDBErr("清除 pending 圖片失敗", err))
		return
	}

	c.Set("data", count)
}
