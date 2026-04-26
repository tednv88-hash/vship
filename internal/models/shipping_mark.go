package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShippingMark represents a user's shipping mark (唛头) for identifying packages
type ShippingMark struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Code        string     `gorm:"type:varchar(100);not null" json:"code"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (sm *ShippingMark) TableName() string {
	return "shipping_marks"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (sm *ShippingMark) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	if sm.ExtraFields == nil {
		sm.ExtraFields = make(JSONB)
	}
	return nil
}
