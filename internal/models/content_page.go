package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContentPage represents a static content page (about, privacy, terms)
type ContentPage struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_content_pages_tenant_slug" json:"tenant_id,omitempty"`
	Slug        string     `gorm:"type:varchar(100);uniqueIndex:idx_content_pages_tenant_slug" json:"slug"`
	Title       string     `gorm:"type:varchar(255)" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *ContentPage) TableName() string {
	return "content_pages"
}

func (x *ContentPage) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
