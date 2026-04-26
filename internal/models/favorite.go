package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Favorite represents a user's favorited goods item
type Favorite struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_favorites_tenant_user_goods" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_favorites_tenant_user_goods" json:"user_id"`
	GoodsID     uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_favorites_tenant_user_goods" json:"goods_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Favorite) TableName() string {
	return "favorites"
}

func (x *Favorite) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
