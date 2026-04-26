package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PageDesign represents a custom page design layout
type PageDesign struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	PageData    JSONB      `gorm:"type:jsonb" json:"page_data"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	IsDefault   bool       `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *PageDesign) TableName() string {
	return "page_designs"
}

func (x *PageDesign) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
