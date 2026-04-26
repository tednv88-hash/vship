package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BatchTemplate represents a template for batch operations
type BatchTemplate struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name            string     `gorm:"type:varchar(255)" json:"name"`
	RouteID         *uuid.UUID `gorm:"type:uuid" json:"route_id,omitempty"`
	Description     string     `gorm:"type:text" json:"description"`
	DefaultSettings JSONB      `gorm:"type:jsonb" json:"default_settings"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TrashedAt       *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields     JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *BatchTemplate) TableName() string {
	return "batch_templates"
}

func (x *BatchTemplate) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
