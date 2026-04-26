package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RechargePlan represents a predefined recharge plan with gift amounts
type RechargePlan struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	GiftAmount  float64    `gorm:"type:numeric(12,2);default:0" json:"gift_amount"`
	GiftPoints  int        `gorm:"default:0" json:"gift_points"`
	SortOrder   int        `json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *RechargePlan) TableName() string {
	return "recharge_plans"
}

func (x *RechargePlan) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
