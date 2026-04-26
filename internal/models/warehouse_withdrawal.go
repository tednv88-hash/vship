package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseWithdrawal represents a withdrawal request from a warehouse
type WarehouseWithdrawal struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID uuid.UUID  `gorm:"type:uuid" json:"warehouse_id"`
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

func (x *WarehouseWithdrawal) TableName() string {
	return "warehouse_withdrawals"
}

func (x *WarehouseWithdrawal) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
