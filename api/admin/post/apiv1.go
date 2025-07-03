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
		apiGroup.GET("/:id", api.GetPostByID)
		apiGroup.POST("", api.CreatePost)
		apiGroup.PATCH("/:id", api.UpdatePost)
		apiGroup.DELETE("/:id", api.DeletePost)
		apiGroup.GET("/category/:id", api.GetPostsByCategory)
		apiGroup.GET("/upload-url", api.GetPresignedUploadURL)
		apiGroup.GET("/about", api.GetAboutMe)
		apiGroup.POST("/about", api.UpdateAboutMe)
	}
}

// 文章列表
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

// 文章詳情
func (api *PostAPI) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	post, err := api.service.GetPostByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", post)
}

// 建立文章
func (api *PostAPI) CreatePost(c *gin.Context) {
	var req CreatePostDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.ErrValidation)
		return
	}
	result, err := api.service.CreatePost(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", result)
}

// 編輯文章
func (api *PostAPI) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePostDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.ErrValidation)
		return
	}
	result, err := api.service.UpdatePost(id, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", result)
}

// 刪除文章
func (api *PostAPI) DeletePost(c *gin.Context) {
	id := c.Param("id")
	err := api.service.DeletePost(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// 取得分類文章
func (api *PostAPI) GetPostsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	posts, err := api.service.GetPostsByCategory(categoryID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", posts)
}

// 提供前端上傳圖片用的預簽名 URL
func (api *PostAPI) GetPresignedUploadURL(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		c.Error(middleware.ErrBadRequest)
		return
	}
	url, err := api.service.GeneratePresignedUploadURL(filename)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", url)
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

// 更新關於我內容
func (api *PostAPI) UpdateAboutMe(c *gin.Context) {
	var req UpdateAboutMeDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.ErrValidation)
		return
	}

	updated, err := api.service.UpdateAboutMe(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", updated)
}
