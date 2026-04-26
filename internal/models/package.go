package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Package represents a shipment package in the consolidation system
type Package struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID         *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	TrackingNumber   string     `gorm:"type:varchar(255);not null;index" json:"tracking_number"`
	ShippingMarkID   *uuid.UUID `gorm:"type:uuid;index" json:"shipping_mark_id,omitempty"`
	UserID           *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	WarehouseID      *uuid.UUID `gorm:"type:uuid" json:"warehouse_id,omitempty"`
	ShelfID          *uuid.UUID `gorm:"type:uuid" json:"shelf_id,omitempty"`
	CategoryID       *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Status           string     `gorm:"type:varchar(50);default:'forecast'" json:"status"`
	Source           string     `gorm:"type:varchar(50)" json:"source,omitempty"`
	Weight           float64    `gorm:"type:decimal(10,3);default:0" json:"weight"`
	Length           float64    `gorm:"type:decimal(10,2);default:0" json:"length"`
	Width            float64    `gorm:"type:decimal(10,2);default:0" json:"width"`
	Height           float64    `gorm:"type:decimal(10,2);default:0" json:"height"`
	VolumeWeight     float64    `gorm:"type:decimal(10,3);default:0" json:"volume_weight"`
	ChargeableWeight float64    `gorm:"type:decimal(10,3);default:0" json:"chargeable_weight"`
	DeclaredValue    float64    `gorm:"type:decimal(12,2);default:0" json:"declared_value"`
	DeclaredCurrency string     `gorm:"type:varchar(10)" json:"declared_currency,omitempty"`
	ItemDescription  string     `gorm:"type:text" json:"item_description,omitempty"`
	Remark           string     `gorm:"type:text" json:"remark,omitempty"`
	IsProblem        bool       `gorm:"default:false" json:"is_problem"`
	ProblemReason    string     `gorm:"type:text" json:"problem_reason,omitempty"`
	IsReturned       bool       `gorm:"default:false" json:"is_returned"`
	IsAppointment    bool       `gorm:"default:false" json:"is_appointment"`
	AppointmentAt    *time.Time `gorm:"type:timestamp with time zone" json:"appointment_at,omitempty"`
	ReceivedAt       *time.Time `gorm:"type:timestamp with time zone" json:"received_at,omitempty"`
	ShelvedAt        *time.Time `gorm:"type:timestamp with time zone" json:"shelved_at,omitempty"`
	InspectedAt      *time.Time `gorm:"type:timestamp with time zone" json:"inspected_at,omitempty"`
	PackedAt         *time.Time `gorm:"type:timestamp with time zone" json:"packed_at,omitempty"`
	ShippedAt        *time.Time `gorm:"type:timestamp with time zone" json:"shipped_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	TrashedAt        *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields      JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (p *Package) TableName() string {
	return "packages"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (p *Package) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.ExtraFields == nil {
		p.ExtraFields = make(JSONB)
	}
	return nil
}
