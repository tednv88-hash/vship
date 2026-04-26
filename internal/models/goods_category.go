package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoodsCategory represents a category for goods
type GoodsCategory struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	ParentID    *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	ImageURL    string     `gorm:"type:varchar(500)" json:"image_url"`
	SortOrder   int        `json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *GoodsCategory) TableName() string {
	return "goods_categories"
}

func (x *GoodsCategory) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
