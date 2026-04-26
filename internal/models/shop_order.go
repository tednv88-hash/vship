package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShopOrder represents a shop order
type ShopOrder struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderNo         string     `gorm:"type:varchar(100);uniqueIndex:idx_shop_orders_tenant_order_no" json:"order_no"`
	UserID          *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Status          string     `gorm:"type:varchar(50);default:'pending_payment'" json:"status"`
	TotalAmount     float64    `gorm:"type:numeric(12,2)" json:"total_amount"`
	DiscountAmount  float64    `gorm:"type:numeric(12,2);default:0" json:"discount_amount"`
	ShippingFee     float64    `gorm:"type:numeric(12,2);default:0" json:"shipping_fee"`
	PayAmount       float64    `gorm:"type:numeric(12,2)" json:"pay_amount"`
	PayMethod       string     `gorm:"type:varchar(50)" json:"pay_method"`
	PayTime         *time.Time `json:"pay_time,omitempty"`
	DeliveryTime    *time.Time `json:"delivery_time,omitempty"`
	ReceiveTime     *time.Time `json:"receive_time,omitempty"`
	Remark          string     `gorm:"type:text" json:"remark"`
	ReceiverName    string     `gorm:"type:varchar(255)" json:"receiver_name"`
	ReceiverPhone   string     `gorm:"type:varchar(50)" json:"receiver_phone"`
	ReceiverAddress string     `gorm:"type:text" json:"receiver_address"`
	ExpressCompany  string     `gorm:"type:varchar(255)" json:"express_company"`
	ExpressNo       string     `gorm:"type:varchar(255)" json:"express_no"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TrashedAt       *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields     JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *ShopOrder) TableName() string {
	return "shop_orders"
}

func (x *ShopOrder) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
