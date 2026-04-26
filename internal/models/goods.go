package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Goods represents a product in the shop
type Goods struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	CategoryID    *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Name          string     `gorm:"type:varchar(255)" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	ImageURL      string     `gorm:"type:varchar(500)" json:"image_url"`
	Images        JSONB      `gorm:"type:jsonb" json:"images"`
	Price         float64    `gorm:"type:numeric(12,2)" json:"price"`
	OriginalPrice float64    `gorm:"type:numeric(12,2)" json:"original_price"`
	Stock         int        `json:"stock"`
	SalesCount    int        `gorm:"default:0" json:"sales_count"`
	Unit          string     `gorm:"type:varchar(50)" json:"unit"`
	Weight        float64    `gorm:"type:numeric(10,3)" json:"weight"`
	SortOrder     int        `json:"sort_order"`
	Status        string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TrashedAt     *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Goods) TableName() string {
	return "goods"
}

func (x *Goods) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
