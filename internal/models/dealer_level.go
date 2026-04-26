package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DealerLevel represents a level in the dealer hierarchy
type DealerLevel struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID         *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name             string     `gorm:"type:varchar(255)" json:"name"`
	LevelNo          int        `json:"level_no"`
	CommissionRate1  float64    `gorm:"type:numeric(5,2)" json:"commission_rate_1"`
	CommissionRate2  float64    `gorm:"type:numeric(5,2)" json:"commission_rate_2"`
	CommissionRate3  float64    `gorm:"type:numeric(5,2)" json:"commission_rate_3"`
	UpgradeCondition JSONB      `gorm:"type:jsonb" json:"upgrade_condition"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	TrashedAt        *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields      JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *DealerLevel) TableName() string {
	return "dealer_levels"
}

func (x *DealerLevel) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
