package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentAudit represents a payment audit record for order payments
type PaymentAudit struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID       *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	UserID         *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Amount         float64    `gorm:"type:numeric(12,2);default:0" json:"amount"`
	Currency       string     `gorm:"type:varchar(10);default:'TWD'" json:"currency"`
	PaymentMethod  string     `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	CertificateURL string     `gorm:"type:varchar(500)" json:"certificate_url,omitempty"`
	Status         string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	ReviewedByID   *uuid.UUID `gorm:"type:uuid" json:"reviewed_by_id,omitempty"`
	ReviewedAt     *time.Time `gorm:"type:timestamp with time zone" json:"reviewed_at,omitempty"`
	RejectReason   string     `gorm:"type:text" json:"reject_reason,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ExtraFields    JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (p *PaymentAudit) TableName() string {
	return "payment_audits"
}

// BeforeCreate sets UUID before creating
func (p *PaymentAudit) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.ExtraFields == nil {
		p.ExtraFields = make(JSONB)
	}
	return nil
}
