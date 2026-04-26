package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RemittanceCertificate represents a remittance certificate uploaded by a user
type RemittanceCertificate struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID         *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID           uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	OrderID          *uuid.UUID `gorm:"type:uuid" json:"order_id,omitempty"`
	Amount           float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Currency         string     `gorm:"type:varchar(10)" json:"currency"`
	BankAccountID    *uuid.UUID `gorm:"type:uuid" json:"bank_account_id,omitempty"`
	CertificateImage string     `gorm:"type:varchar(500)" json:"certificate_image"`
	Remark           string     `gorm:"type:text" json:"remark"`
	Status           string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	AuditRemark      string     `gorm:"type:text" json:"audit_remark"`
	AuditorID        *uuid.UUID `gorm:"type:uuid" json:"auditor_id,omitempty"`
	AuditedAt        *time.Time `json:"audited_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	TrashedAt        *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields      JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *RemittanceCertificate) TableName() string {
	return "remittance_certificates"
}

func (x *RemittanceCertificate) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
