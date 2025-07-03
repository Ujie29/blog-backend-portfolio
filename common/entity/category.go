package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID          uint      `bun:",pk,autoincrement,notnull"` // 主鍵，不為 null
	Name        string    `bun:",notnull"`                  // 分類名稱，不為 null
	Parent      *uint     // 上層分類，不為 null
	Slug        string    `bun:",unique,notnull"`                    // 分類的 URL slug
	CreatedAt   time.Time `bun:",notnull,default:current_timestamp"` // 建立時間，不為 null
	UpdatedAt   time.Time `bun:",notnull,default:current_timestamp"` // 更新時間，不為 null
	HasChildren bool      `bun:"has_children,notnull"`               // 是否有子分類
	SortOrder   int       `bun:"sort_order,notnull,default:0"`       // 排序用欄位
}
