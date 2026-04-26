package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DealerOrder represents a commission record for a dealer on an order
type DealerOrder struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID         *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	DealerID         uuid.UUID  `gorm:"type:uuid" json:"dealer_id"`
	OrderID          uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	OrderType        string     `gorm:"type:varchar(50)" json:"order_type"`
	CommissionAmount float64    `gorm:"type:numeric(12,2)" json:"commission_amount"`
	Status           string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	SettledAt        *time.Time `json:"settled_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	ExtraFields      JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *DealerOrder) TableName() string {
	return "dealer_orders"
}

func (x *DealerOrder) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
