package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CartItem represents a shopping cart item
type CartItem struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	GoodsID     uuid.UUID  `gorm:"type:uuid" json:"goods_id"`
	SkuID       *uuid.UUID `gorm:"type:uuid" json:"sku_id,omitempty"`
	Quantity    int        `gorm:"default:1" json:"quantity"`
	Selected    bool       `gorm:"default:true" json:"selected"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *CartItem) TableName() string {
	return "cart_items"
}

func (x *CartItem) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
