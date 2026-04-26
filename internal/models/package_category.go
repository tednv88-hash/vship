package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageCategory represents a category for classifying packages
type PackageCategory struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (p *PackageCategory) TableName() string {
	return "package_categories"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (p *PackageCategory) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.ExtraFields == nil {
		p.ExtraFields = make(JSONB)
	}
	return nil
}
