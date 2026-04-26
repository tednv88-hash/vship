package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoutePricingTier represents a pricing tier for a shipping route
type RoutePricingTier struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID              *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	RouteID               uuid.UUID  `gorm:"type:uuid;not null;index" json:"route_id"`
	MemberLevelID         *uuid.UUID `gorm:"type:uuid" json:"member_level_id,omitempty"`
	WeightMin             float64    `gorm:"type:numeric(10,2);default:0" json:"weight_min"`
	WeightMax             float64    `gorm:"type:numeric(10,2);default:0" json:"weight_max"`
	UnitPrice             float64    `gorm:"type:numeric(10,2);default:0" json:"unit_price"`
	FirstWeight           float64    `gorm:"type:numeric(10,2);default:0" json:"first_weight"`
	FirstWeightPrice      float64    `gorm:"type:numeric(10,2);default:0" json:"first_weight_price"`
	AdditionalWeightPrice float64    `gorm:"type:numeric(10,2);default:0" json:"additional_weight_price"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (r *RoutePricingTier) TableName() string {
	return "route_pricing_tiers"
}

// BeforeCreate sets UUID before creating
func (r *RoutePricingTier) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
