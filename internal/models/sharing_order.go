package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SharingOrder represents a commission record from sharing
type SharingOrder struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID         *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID           uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	SharerID         uuid.UUID  `gorm:"type:uuid" json:"sharer_id"`
	OrderID          uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	OrderType        string     `gorm:"type:varchar(50)" json:"order_type"`
	CommissionAmount float64    `gorm:"type:numeric(12,2)" json:"commission_amount"`
	Status           string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	ExtraFields      JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *SharingOrder) TableName() string {
	return "sharing_orders"
}

func (x *SharingOrder) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
