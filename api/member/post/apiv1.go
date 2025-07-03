package post

import (
	"blog-backend/common/middleware"

	"github.com/gin-gonic/gin"
)

type PostAPI struct {
	service PostService
}

func NewPostAPI(service PostService) *PostAPI {
	return &PostAPI{
		service: service,
	}
}

func (api *PostAPI) RegisterRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api/post")
	{
		apiGroup.GET("", api.GetPostList)
		apiGroup.GET("/:slug", api.GetPostBySlug)
		apiGroup.GET("/category/:slug", api.GetPostsByCategory)
		apiGroup.GET("/about", api.GetAboutMe)
		apiGroup.POST("/randomCategoryPost", api.GetRandomPostsByCategory)
	}
}

// 取得所有文章
func (api *PostAPI) GetPostList(c *gin.Context) {
	var req GetPostListDto
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(middleware.ErrBadRequest)
		return
	}
	result, err := api.service.GetPostList(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", result)
}

// 取得單一文章
func (api *PostAPI) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	post, err := api.service.GetPostBySlug(slug)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", post)
}

// 取得分類文章
func (api *PostAPI) GetPostsByCategory(c *gin.Context) {
	slug := c.Param("slug")

	var req GetPostListDto
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(middleware.ErrBadRequest)
		return
	}
	result, err := api.service.GetPostsByCategory(slug, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", result)
}

// 取得關於我內容
func (api *PostAPI) GetAboutMe(c *gin.Context) {
	about, err := api.service.GetAboutMe()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", about)
}

// 取得分類底下隨機六篇文章
func (api *PostAPI) GetRandomPostsByCategory(c *gin.Context) {
	var dto GetRandomPostsByCategoryDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(middleware.ErrBadRequest)
		return
	}

	posts, err := api.service.GetRandomPostsByCategory(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", posts)
}
