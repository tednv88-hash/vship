package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RouteCategory represents a category for classifying shipping routes
type RouteCategory struct {
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

func (r *RouteCategory) TableName() string {
	return "route_categories"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (r *RouteCategory) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.ExtraFields == nil {
		r.ExtraFields = make(JSONB)
	}
	return nil
}
