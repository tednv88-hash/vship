package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoodsReview represents a review for a goods item
type GoodsReview struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	GoodsID     uuid.UUID  `gorm:"type:uuid" json:"goods_id"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	OrderID     *uuid.UUID `gorm:"type:uuid" json:"order_id,omitempty"`
	Content     string     `gorm:"type:text" json:"content"`
	Rating      int        `json:"rating"`
	Images      JSONB      `gorm:"type:jsonb" json:"images"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Reply       string     `gorm:"type:text" json:"reply"`
	RepliedAt   *time.Time `json:"replied_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *GoodsReview) TableName() string {
	return "goods_reviews"
}

func (x *GoodsReview) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
