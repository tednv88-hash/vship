package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TrackingTemplate represents a common tracking status template (常用軌跡)
type TrackingTemplate struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (t *TrackingTemplate) TableName() string {
	return "tracking_templates"
}

// BeforeCreate sets UUID before creating
func (t *TrackingTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.ExtraFields == nil {
		t.ExtraFields = make(JSONB)
	}
	return nil
}
