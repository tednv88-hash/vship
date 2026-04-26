package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberLevel represents a membership tier with discount rates
type MemberLevel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Code         string     `gorm:"type:varchar(50)" json:"code"`
	DiscountRate float64    `gorm:"type:decimal(5,2);default:0" json:"discount_rate"`
	MinSpend     float64    `gorm:"type:decimal(12,2);default:0" json:"min_spend"`
	SortOrder    int        `gorm:"default:0" json:"sort_order"`
	Status       string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TrashedAt    *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields  JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (m *MemberLevel) TableName() string {
	return "member_levels"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (m *MemberLevel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.ExtraFields == nil {
		m.ExtraFields = make(JSONB)
	}
	return nil
}
