package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserDiscount represents a discount rate assigned to a user for a route
type UserDiscount struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID       uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	RouteID      *uuid.UUID `gorm:"type:uuid" json:"route_id,omitempty"`
	DiscountRate float64    `gorm:"type:numeric(5,2)" json:"discount_rate"`
	Remark       string     `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TrashedAt    *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields  JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *UserDiscount) TableName() string {
	return "user_discounts"
}

func (x *UserDiscount) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
