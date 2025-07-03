package model

import "blog-backend/common/utils"

type PaginatedResponse[T any] struct {
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalCount int                  `json:"totalCount"`
	Data       utils.NonNilSlice[T] `json:"data"`
}
