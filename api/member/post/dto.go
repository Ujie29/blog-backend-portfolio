package post

import "time"

type GetPostListDto struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type PostListDto struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	CoverImageUrl string    `json:"coverImageUrl"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PostDto struct {
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	CategoryID    uint      `json:"categoryId"`
	CoverImageUrl string    `json:"coverImageUrl"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AboutMeDto struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetRandomPostsByCategoryDto struct {
	CategoryID uint   `json:"categoryId"`
	Slug       string `json:"slug"`
}
