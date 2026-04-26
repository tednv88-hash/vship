package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConsolidationOrder represents a consolidation shipping order (集運訂單)
type ConsolidationOrder struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID          *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderNumber       string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_order_number" json:"order_number"`
	UserID            *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ShippingRouteID   *uuid.UUID `gorm:"type:uuid" json:"shipping_route_id,omitempty"`
	WarehouseID       *uuid.UUID `gorm:"type:uuid" json:"warehouse_id,omitempty"`
	MemberLevelID     *uuid.UUID `gorm:"type:uuid" json:"member_level_id,omitempty"`
	Status            string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	TotalWeight       float64    `gorm:"type:numeric(10,3);default:0" json:"total_weight"`
	TotalVolumeWeight float64    `gorm:"type:numeric(10,3);default:0" json:"total_volume_weight"`
	ChargeableWeight  float64    `gorm:"type:numeric(10,3);default:0" json:"chargeable_weight"`
	ShippingFee       float64    `gorm:"type:numeric(12,2);default:0" json:"shipping_fee"`
	InsuranceFee      float64    `gorm:"type:numeric(12,2);default:0" json:"insurance_fee"`
	ServiceFee        float64    `gorm:"type:numeric(12,2);default:0" json:"service_fee"`
	ConsumableFee     float64    `gorm:"type:numeric(12,2);default:0" json:"consumable_fee"`
	TotalAmount       float64    `gorm:"type:numeric(12,2);default:0" json:"total_amount"`
	PaidAmount        float64    `gorm:"type:numeric(12,2);default:0" json:"paid_amount"`
	Currency          string     `gorm:"type:varchar(10);default:'TWD'" json:"currency"`
	PaymentMethod     string     `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentStatus     string     `gorm:"type:varchar(50);default:'unpaid'" json:"payment_status"`
	RecipientName     string     `gorm:"type:varchar(255)" json:"recipient_name,omitempty"`
	RecipientPhone    string     `gorm:"type:varchar(50)" json:"recipient_phone,omitempty"`
	RecipientAddress  string     `gorm:"type:text" json:"recipient_address,omitempty"`
	RecipientCity     string     `gorm:"type:varchar(100)" json:"recipient_city,omitempty"`
	RecipientState    string     `gorm:"type:varchar(100)" json:"recipient_state,omitempty"`
	RecipientZip      string     `gorm:"type:varchar(20)" json:"recipient_zip,omitempty"`
	RecipientCountry  string     `gorm:"type:varchar(100)" json:"recipient_country,omitempty"`
	Remark            string     `gorm:"type:text" json:"remark,omitempty"`
	PaidAt            *time.Time `gorm:"type:timestamp with time zone" json:"paid_at,omitempty"`
	PackedAt          *time.Time `gorm:"type:timestamp with time zone" json:"packed_at,omitempty"`
	ShippedAt         *time.Time `gorm:"type:timestamp with time zone" json:"shipped_at,omitempty"`
	ArrivedAt         *time.Time `gorm:"type:timestamp with time zone" json:"arrived_at,omitempty"`
	CompletedAt       *time.Time `gorm:"type:timestamp with time zone" json:"completed_at,omitempty"`
	CancelledAt       *time.Time `gorm:"type:timestamp with time zone" json:"cancelled_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	TrashedAt         *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields       JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (c *ConsolidationOrder) TableName() string {
	return "consolidation_orders"
}

// BeforeCreate sets UUID before creating
func (c *ConsolidationOrder) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExtraFields == nil {
		c.ExtraFields = make(JSONB)
	}
	return nil
}
