package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CouponReceiveLog represents a log of coupon received by a user
type CouponReceiveLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	CouponID    uuid.UUID  `gorm:"type:uuid" json:"coupon_id"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	Status      string     `gorm:"type:varchar(50);default:'unused'" json:"status"`
	UsedAt      *time.Time `json:"used_at,omitempty"`
	OrderID     *uuid.UUID `gorm:"type:uuid" json:"order_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *CouponReceiveLog) TableName() string {
	return "coupon_receive_logs"
}

func (x *CouponReceiveLog) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
