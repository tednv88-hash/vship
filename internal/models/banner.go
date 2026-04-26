package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Banner represents a promotional banner
type Banner struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	ImageURL    string     `gorm:"type:varchar(500);not null" json:"image_url"`
	LinkURL     string     `gorm:"type:varchar(500)" json:"link_url,omitempty"`
	Position    string     `gorm:"type:varchar(100)" json:"position,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	StartAt     *time.Time `gorm:"type:timestamp with time zone" json:"start_at,omitempty"`
	EndAt       *time.Time `gorm:"type:timestamp with time zone" json:"end_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (b *Banner) TableName() string {
	return "banners"
}

// BeforeCreate sets UUID before creating
func (b *Banner) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.ExtraFields == nil {
		b.ExtraFields = make(JSONB)
	}
	return nil
}
