package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminRole represents an administrative role with permissions
type AdminRole struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Permissions JSONB      `gorm:"type:jsonb" json:"permissions"`
	SortOrder   int        `json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *AdminRole) TableName() string {
	return "admin_roles"
}

func (x *AdminRole) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
