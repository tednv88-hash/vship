package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InsuranceProduct represents an insurance product for shipments
type InsuranceProduct struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	PremiumRate float64    `gorm:"type:numeric(8,4);default:0" json:"premium_rate"`
	MinPremium  float64    `gorm:"type:numeric(10,2);default:0" json:"min_premium"`
	MaxCoverage float64    `gorm:"type:numeric(12,2);default:0" json:"max_coverage"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (i *InsuranceProduct) TableName() string {
	return "insurance_products"
}

// BeforeCreate sets UUID before creating
func (i *InsuranceProduct) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.ExtraFields == nil {
		i.ExtraFields = make(JSONB)
	}
	return nil
}
