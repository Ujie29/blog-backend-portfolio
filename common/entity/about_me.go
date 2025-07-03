package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type AboutMe struct {
	bun.BaseModel `bun:"table:about_me"`

	ID          uint      `bun:"id,pk,autoincrement"`
	HtmlContent string    `bun:"html_content,type:text,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}
