package entity

// import (
// 	"github.com/uptrace/bun"
// )

// type ContentBlock struct {
// 	bun.BaseModel `bun:"table:contentblock"`

// 	ID          uint   `bun:",pk,autoincrement,notnull"`
// 	HTMLContent string `bun:",notnull"`
// 	PostID      uint   `bun:",notnull"`

// 	Post *Post `bun:"rel:belongs-to,join:post_id=id"` // 關聯回 post
// }
