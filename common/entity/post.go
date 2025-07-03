package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:posts"`

	ID            uint      `bun:",pk,autoincrement,notnull"`          // 主鍵，自動遞增，不為 null
	Title         string    `bun:",notnull"`                           // 標題，不為 null
	CategoryID    uint      `bun:",notnull"`                           // 分類 ID，不為 null
	IsPublished   bool      `bun:",notnull"`                           // 是否發佈，不為 null
	Slug          string    `bun:",unique,notnull"`                    // 對 SEO 友善的唯一識別 slug
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp"` // 建立時間，不為 null
	UpdatedAt     time.Time `bun:",notnull,default:current_timestamp"` // 更新時間，不為 null
	CoverImageUrl string    `bun:",notnull"`                           // 封面圖片，不為 null
	Content       string    `bun:"content"`                            // 文章內容
	IsDeleted     bool      `bun:",notnull,default:false"`             // 軟刪除欄位
	NeedsRefresh  bool      `bun:",notnull,default:false"`             // 內容變更時觸發刷新
}
