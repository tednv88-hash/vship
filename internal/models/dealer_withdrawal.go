package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DealerWithdrawal represents a withdrawal request from a dealer
type DealerWithdrawal struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	DealerID    uuid.UUID  `gorm:"type:uuid" json:"dealer_id"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Method      string     `gorm:"type:varchar(50)" json:"method"`
	AccountInfo JSONB      `gorm:"type:jsonb" json:"account_info"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	AuditRemark string     `gorm:"type:text" json:"audit_remark"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *DealerWithdrawal) TableName() string {
	return "dealer_withdrawals"
}

func (x *DealerWithdrawal) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
