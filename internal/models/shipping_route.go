package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShippingRoute represents a consolidation shipping route (集運線路)
type ShippingRoute struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID               *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name                   string     `gorm:"type:varchar(255);not null" json:"name"`
	OriginCountryID        *uuid.UUID `gorm:"type:uuid" json:"origin_country_id,omitempty"`
	DestinationCountryID   *uuid.UUID `gorm:"type:uuid" json:"destination_country_id,omitempty"`
	OriginWarehouseID      *uuid.UUID `gorm:"type:uuid" json:"origin_warehouse_id,omitempty"`
	DestinationWarehouseID *uuid.UUID `gorm:"type:uuid" json:"destination_warehouse_id,omitempty"`
	CategoryID             *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	TransportMode          string     `gorm:"type:varchar(50);not null" json:"transport_mode"`
	BillingMode            string     `gorm:"type:varchar(50);not null" json:"billing_mode"`
	WeightUnit             string     `gorm:"type:varchar(10);default:'KG'" json:"weight_unit"`
	VolumeWeightRatio      float64    `gorm:"type:numeric(10,2);default:5000" json:"volume_weight_ratio"`
	RoundingRule           string     `gorm:"type:varchar(50)" json:"rounding_rule,omitempty"`
	MultiBoxMode           string     `gorm:"type:varchar(50);default:'separate'" json:"multi_box_mode"`
	EstimatedDays          int        `gorm:"default:0" json:"estimated_days"`
	Status                 string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder              int        `gorm:"default:0" json:"sort_order"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	TrashedAt              *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields            JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (s *ShippingRoute) TableName() string {
	return "shipping_routes"
}

// BeforeCreate sets UUID before creating
func (s *ShippingRoute) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.ExtraFields == nil {
		s.ExtraFields = make(JSONB)
	}
	return nil
}
