package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShopOrderItem represents an item within a shop order
type ShopOrderItem struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID     uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	GoodsID     uuid.UUID  `gorm:"type:uuid" json:"goods_id"`
	GoodsName   string     `gorm:"type:varchar(255)" json:"goods_name"`
	GoodsImage  string     `gorm:"type:varchar(500)" json:"goods_image"`
	Price       float64    `gorm:"type:numeric(12,2)" json:"price"`
	Quantity    int        `json:"quantity"`
	TotalAmount float64    `gorm:"type:numeric(12,2)" json:"total_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *ShopOrderItem) TableName() string {
	return "shop_order_items"
}

func (x *ShopOrderItem) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
