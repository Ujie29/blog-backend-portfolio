package post

import "time"

type GetPostListDto struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type PostListDto struct {
	SortID    int       `json:"sortId"`
	Id        uint      `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostDto struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Summary       string    `json:"summary"`
	CoverImageUrl string    `json:"coverImageUrl"`
	IsPublished   bool      `json:"isPublished"`
	CategoryID    uint      `json:"categoryId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Slug          string    `json:"slug"`
}

type UpdatePostDto = CreatePostDto

type CreatePostDto struct {
	Title         string `json:"title" binding:"required"`
	CoverImageUrl string `json:"coverImageUrl"`
	Content       string `json:"content" binding:"required"`
	CategoryID    uint   `json:"categoryId"`
	IsPublished   bool   `json:"isPublished"`
	Slug          string `json:"slug" binding:"required"`
}

type UploadUrlDto struct {
	UploadUrl string `json:"uploadUrl"`
	ImageUrl  string `json:"imageUrl"`
	ExpiresIn int64  `json:"expiresIn"`
	Filename  string `json:"filename"`
}

type UpdateAboutMeDto struct {
	Content string `json:"content" binding:"required"`
}

type AboutMeDto struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updatedAt"`
}
