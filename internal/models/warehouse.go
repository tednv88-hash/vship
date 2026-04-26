package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Warehouse represents a physical warehouse location
type Warehouse struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
	Code          string     `gorm:"type:varchar(50)" json:"code"`
	CountryID     *uuid.UUID `gorm:"type:uuid" json:"country_id,omitempty"`
	Address       string     `gorm:"type:text" json:"address,omitempty"`
	Phone         string     `gorm:"type:varchar(50)" json:"phone,omitempty"`
	ContactPerson string     `gorm:"type:varchar(255)" json:"contact_person,omitempty"`
	Status        string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	IsDefault     bool       `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TrashedAt     *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (w *Warehouse) TableName() string {
	return "warehouses"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (w *Warehouse) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.ExtraFields == nil {
		w.ExtraFields = make(JSONB)
	}
	return nil
}
