package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dealer represents a dealer in the distribution system
type Dealer struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID            *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID              uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	LevelID             *uuid.UUID `gorm:"type:uuid" json:"level_id,omitempty"`
	ParentID            *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	CommissionRate      float64    `gorm:"type:numeric(5,2)" json:"commission_rate"`
	TotalCommission     float64    `gorm:"type:numeric(12,2);default:0" json:"total_commission"`
	WithdrawnCommission float64    `gorm:"type:numeric(12,2);default:0" json:"withdrawn_commission"`
	Status              string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	TrashedAt           *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields         JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Dealer) TableName() string {
	return "dealers"
}

func (x *Dealer) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
