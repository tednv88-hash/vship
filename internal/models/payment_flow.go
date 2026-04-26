package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentFlow represents a payment transaction flow record
type PaymentFlow struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderType   string     `gorm:"type:varchar(50)" json:"order_type"`
	OrderID     uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	OrderNo     string     `gorm:"type:varchar(100)" json:"order_no"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	PayMethod   string     `gorm:"type:varchar(50)" json:"pay_method"`
	PayNo       string     `gorm:"type:varchar(255)" json:"pay_no"`
	Status      string     `gorm:"type:varchar(50)" json:"status"`
	PayTime     *time.Time `json:"pay_time,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *PaymentFlow) TableName() string {
	return "payment_flows"
}

func (x *PaymentFlow) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
