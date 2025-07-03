package category

import (
	"context"
	"database/sql"
	"errors"

	"blog-backend/common/entity"
	"blog-backend/common/middleware"

	"github.com/uptrace/bun"
)

type CategoryService interface {
	GetCategoryTree() ([]*CategoryResponse, error)
	GetCategoryByID(id string) (entity.Category, error)
}

type categoryServiceImpl struct {
	db *bun.DB
}

func NewCategoryService(db *bun.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

func (s *categoryServiceImpl) GetCategoryTree() ([]*CategoryResponse, error) {
	var categories []entity.Category
	err := s.db.NewSelect().
		Model(&categories).
		Scan(context.Background())
	if err != nil {
		return nil, middleware.WrapDBErr("查詢分類失敗", err)
	}

	idToNode := make(map[uint]*CategoryResponse)
	childrenMap := make(map[uint][]*CategoryResponse)

	// 先建立所有節點
	for _, cat := range categories {
		idToNode[cat.ID] = &CategoryResponse{
			ID:       cat.ID,
			Name:     cat.Name,
			Children: []*CategoryResponse{},
		}
	}

	// 建立 parentID -> 子節點 map
	for _, cat := range categories {
		if cat.Parent != nil {
			childrenMap[*cat.Parent] = append(childrenMap[*cat.Parent], idToNode[cat.ID])
		}
	}

	// 遞迴掛上子節點
	var attachChildren func(node *CategoryResponse)
	attachChildren = func(node *CategoryResponse) {
		if children, ok := childrenMap[node.ID]; ok {
			node.Children = children
		} else {
			node.Children = []*CategoryResponse{} // 👈 確保不是 nil
		}
		for _, child := range node.Children {
			attachChildren(child)
		}
	}

	var roots []*CategoryResponse
	for _, cat := range categories {
		if cat.Parent == nil {
			root := idToNode[cat.ID]
			attachChildren(root)
			roots = append(roots, root)
		}
	}

	return roots, nil
}

func (s *categoryServiceImpl) GetCategoryByID(id string) (entity.Category, error) {
	var category entity.Category
	err := s.db.NewSelect().
		Model(&category).
		Where("category.id = ?", id).
		Limit(1).
		Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Category{}, middleware.ErrNotFound
	} else if err != nil {
		return entity.Category{}, middleware.WrapDBErr("查詢分類失敗", err)
	}

	return category, nil
}
