package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Currency represents a currency with exchange rate info
type Currency struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Code         string     `gorm:"type:varchar(10);not null" json:"code"`
	Symbol       string     `gorm:"type:varchar(10)" json:"symbol"`
	ExchangeRate float64    `gorm:"type:decimal(18,6);default:1.0" json:"exchange_rate"`
	IsDefault    bool       `gorm:"default:false" json:"is_default"`
	Status       string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TrashedAt    *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields  JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (c *Currency) TableName() string {
	return "currencies"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (c *Currency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExtraFields == nil {
		c.ExtraFields = make(JSONB)
	}
	return nil
}
