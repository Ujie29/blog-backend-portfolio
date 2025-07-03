package category

type CreateCategoryRequest struct {
	Name   string `json:"name" binding:"required"` // 分類名稱，必填
	Parent *uint  `json:"parent"`                  // 可選的上層分類 ID
}

type CategoryResponse struct {
	ID       uint                `json:"id"`       // 分類 ID
	Name     string              `json:"name"`     // 分類名稱
	Children []*CategoryResponse `json:"children"` // 子分類
}
