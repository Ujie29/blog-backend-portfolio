package category

import (
	"github.com/gin-gonic/gin"
)

// CategoryAPI 是 API 層，持有 service
type CategoryAPI struct {
	service CategoryService
}

// NewCategoryAPI 建立新的 CategoryAPI 實例
func NewCategoryAPI(service CategoryService) *CategoryAPI {
	return &CategoryAPI{
		service: service,
	}
}

// RegisterRoutes 註冊路由到 gin Engine
func (api *CategoryAPI) RegisterRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api/category")
	{
		apiGroup.GET("", api.GetCategories)
		apiGroup.GET("/:id", api.GetCategoryByID)
	}
}

// GetCategories 取得分類樹
func (api *CategoryAPI) GetCategories(c *gin.Context) {
	categories, err := api.service.GetCategoryTree()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", categories)
}

// GetCategoryByID 取得單一分類
func (api *CategoryAPI) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := api.service.GetCategoryByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", category)
}
