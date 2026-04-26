package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RechargeOrder represents a balance recharge order
type RechargeOrder struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	OrderNo     string     `gorm:"type:varchar(100)" json:"order_no"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	GiftAmount  float64    `gorm:"type:numeric(12,2);default:0" json:"gift_amount"`
	PayMethod   string     `gorm:"type:varchar(50)" json:"pay_method"`
	PayTime     *time.Time `json:"pay_time,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *RechargeOrder) TableName() string {
	return "recharge_orders"
}

func (x *RechargeOrder) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
