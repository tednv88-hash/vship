package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Country represents a country with ISO codes
type Country struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	NameEN      string     `gorm:"type:varchar(255)" json:"name_en"`
	Code        string     `gorm:"type:varchar(2);not null" json:"code"`
	PhoneCode   string     `gorm:"type:varchar(10)" json:"phone_code"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (c *Country) TableName() string {
	return "countries"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (c *Country) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExtraFields == nil {
		c.ExtraFields = make(JSONB)
	}
	return nil
}
