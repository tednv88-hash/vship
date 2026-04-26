package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShippingBatch represents a shipping batch (批次)
type ShippingBatch struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID               *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	BatchNumber            string     `gorm:"type:varchar(100);not null" json:"batch_number"`
	Name                   string     `gorm:"type:varchar(255)" json:"name,omitempty"`
	Type                   string     `gorm:"type:varchar(50);not null" json:"type"`
	ContainerCode          string     `gorm:"type:varchar(100)" json:"container_code,omitempty"`
	MasterTrackingNumber   string     `gorm:"type:varchar(255)" json:"master_tracking_number,omitempty"`
	OriginWarehouseID      *uuid.UUID `gorm:"type:uuid" json:"origin_warehouse_id,omitempty"`
	DestinationWarehouseID *uuid.UUID `gorm:"type:uuid" json:"destination_warehouse_id,omitempty"`
	LogisticsCompanyID     *uuid.UUID `gorm:"type:uuid" json:"logistics_company_id,omitempty"`
	Status                 string     `gorm:"type:varchar(50);default:'preparing'" json:"status"`
	TotalWeight            float64    `gorm:"type:numeric(12,3);default:0" json:"total_weight"`
	TotalVolume            float64    `gorm:"type:numeric(12,3);default:0" json:"total_volume"`
	TotalOrders            int        `gorm:"default:0" json:"total_orders"`
	DepartedAt             *time.Time `gorm:"type:timestamp with time zone" json:"departed_at,omitempty"`
	ArrivedAt              *time.Time `gorm:"type:timestamp with time zone" json:"arrived_at,omitempty"`
	CompletedAt            *time.Time `gorm:"type:timestamp with time zone" json:"completed_at,omitempty"`
	Remark                 string     `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	TrashedAt              *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields            JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (s *ShippingBatch) TableName() string {
	return "shipping_batches"
}

// BeforeCreate sets UUID before creating
func (s *ShippingBatch) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.ExtraFields == nil {
		s.ExtraFields = make(JSONB)
	}
	return nil
}
