package category

import "time"

type CreateCategoryDto struct {
	Name   string `json:"name" binding:"required"` // 分類名稱，必填
	Slug   string `json:"slug"`
	Parent *uint  `json:"parent"` // 可選的上層分類 ID
}

type CategoryDto struct {
	ID       uint           `json:"id"`   // 分類 ID
	Name     string         `json:"name"` // 分類名稱
	Slug     string         `json:"slug"`
	Children []*CategoryDto `json:"children"` // 子分類
}

type PostListDto struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	CoverImageUrl string    `json:"coverImageUrl"`
	CreatedAt     time.Time `json:"createdAt"`
}
