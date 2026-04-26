package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BrowsingHistory represents a user's browsing history record
type BrowsingHistory struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	GoodsID     uuid.UUID  `gorm:"type:uuid" json:"goods_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *BrowsingHistory) TableName() string {
	return "browsing_histories"
}

func (x *BrowsingHistory) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
