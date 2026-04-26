package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseClerkReview represents a review for a warehouse clerk
type WarehouseClerkReview struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ClerkID     uuid.UUID  `gorm:"type:uuid" json:"clerk_id"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	OrderID     *uuid.UUID `gorm:"type:uuid" json:"order_id,omitempty"`
	Rating      int        `json:"rating"`
	Content     string     `gorm:"type:text" json:"content"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseClerkReview) TableName() string {
	return "warehouse_clerk_reviews"
}

func (x *WarehouseClerkReview) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
