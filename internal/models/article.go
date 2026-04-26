package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Article represents a content article
type Article struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content,omitempty"`
	Category    string     `gorm:"type:varchar(100)" json:"category,omitempty"`
	CoverImage  string     `gorm:"type:varchar(500)" json:"cover_image,omitempty"`
	Author      string     `gorm:"type:varchar(255)" json:"author,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	PublishedAt *time.Time `gorm:"type:timestamp with time zone" json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (a *Article) TableName() string {
	return "articles"
}

// BeforeCreate sets UUID before creating
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.ExtraFields == nil {
		a.ExtraFields = make(JSONB)
	}
	return nil
}
