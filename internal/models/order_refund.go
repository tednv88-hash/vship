package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderRefund represents a refund request for an order
type OrderRefund struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID     uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	OrderItemID *uuid.UUID `gorm:"type:uuid" json:"order_item_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	RefundNo    string     `gorm:"type:varchar(100)" json:"refund_no"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Reason      string     `gorm:"type:varchar(500)" json:"reason"`
	Description string     `gorm:"type:text" json:"description"`
	Images      JSONB      `gorm:"type:jsonb" json:"images"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	AuditRemark string     `gorm:"type:text" json:"audit_remark"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *OrderRefund) TableName() string {
	return "order_refunds"
}

func (x *OrderRefund) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
