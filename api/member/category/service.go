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
	GetCategoryTree() ([]*CategoryDto, error)
	GetCategoryBySlug(slug string) (entity.Category, error)
}

type categoryServiceImpl struct {
	db *bun.DB
}

func NewCategoryService(db *bun.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

func (s *categoryServiceImpl) GetCategoryTree() ([]*CategoryDto, error) {
	var categories []entity.Category
	err := s.db.NewSelect().
		Model(&categories).
		Order("sort_order ASC").
		Scan(context.Background())
	if err != nil {
		return nil, middleware.WrapDBErr("查詢分類樹失敗", err)
	}

	idToNode := make(map[uint]*CategoryDto)
	childrenMap := make(map[uint][]*CategoryDto)

	// 先建立所有節點
	for _, cat := range categories {
		idToNode[cat.ID] = &CategoryDto{
			ID:       cat.ID,
			Name:     cat.Name,
			Slug:     cat.Slug,
			Children: []*CategoryDto{},
		}
	}

	// 建立 parentID -> 子節點 map
	for _, cat := range categories {
		if cat.Parent != nil {
			childrenMap[*cat.Parent] = append(childrenMap[*cat.Parent], idToNode[cat.ID])
		}
	}

	// 遞迴掛上子節點
	var attachChildren func(node *CategoryDto)
	attachChildren = func(node *CategoryDto) {
		if children, ok := childrenMap[node.ID]; ok {
			node.Children = children
		} else {
			node.Children = []*CategoryDto{}
		}
		for _, child := range node.Children {
			attachChildren(child)
		}
	}

	var roots []*CategoryDto
	for _, cat := range categories {
		if cat.Parent == nil {
			root := idToNode[cat.ID]
			attachChildren(root)
			roots = append(roots, root)
		}
	}

	return roots, nil
}

func (s *categoryServiceImpl) GetCategoryBySlug(slug string) (entity.Category, error) {
	var category entity.Category
	err := s.db.NewSelect().
		Model(&category).
		Where("category.slug = ?", slug).
		Limit(1).
		Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Category{}, middleware.ErrNotFound
	} else if err != nil {
		return entity.Category{}, middleware.WrapDBErr("查詢分類失敗", err)
	}

	return category, nil
}
