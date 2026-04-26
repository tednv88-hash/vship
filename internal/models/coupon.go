package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID       *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Code           string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_coupon_code" json:"code"`
	Name           string     `gorm:"type:varchar(255);not null" json:"name"`
	Type           string     `gorm:"type:varchar(50);not null" json:"type"`
	Value          float64    `gorm:"type:numeric(10,2);default:0" json:"value"`
	MinOrderAmount float64    `gorm:"type:numeric(12,2);default:0" json:"min_order_amount"`
	MaxDiscount    float64    `gorm:"type:numeric(12,2);default:0" json:"max_discount"`
	TotalCount     int        `gorm:"default:0" json:"total_count"`
	UsedCount      int        `gorm:"default:0" json:"used_count"`
	MemberLevelID  *uuid.UUID `gorm:"type:uuid" json:"member_level_id,omitempty"`
	Status         string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	StartAt        *time.Time `gorm:"type:timestamp with time zone" json:"start_at,omitempty"`
	EndAt          *time.Time `gorm:"type:timestamp with time zone" json:"end_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	TrashedAt      *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields    JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (c *Coupon) TableName() string {
	return "coupons"
}

// BeforeCreate sets UUID before creating
func (c *Coupon) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExtraFields == nil {
		c.ExtraFields = make(JSONB)
	}
	return nil
}
