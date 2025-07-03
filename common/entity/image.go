package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Image struct {
	bun.BaseModel `bun:"table:images"`

	ID        uuid.UUID  `bun:",pk,type:uuid,default:gen_random_uuid()"`
	URL       string     `bun:",notnull"`
	PostID    uint       `bun:",notnull"`
	Type      string     `bun:",notnull"` // 'inline' 或 'cover'
	Status    string     `bun:",notnull"` // 'active' 或 'pending_delete'
	IsDeleted bool       `bun:"is_deleted" json:"isDeleted"`
	DeletedAt *time.Time `bun:"deleted_at" json:"deletedAt,omitempty"`
	CreatedAt time.Time  `bun:",default:now()"`
	UpdatedAt time.Time  `bun:",default:now()"`
}
