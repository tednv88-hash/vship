package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseApplication represents an application to register a warehouse
type WarehouseApplication struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ApplicantName string     `gorm:"type:varchar(255)" json:"applicant_name"`
	Phone         string     `gorm:"type:varchar(50)" json:"phone"`
	WarehouseName string     `gorm:"type:varchar(255)" json:"warehouse_name"`
	Address       string     `gorm:"type:text" json:"address"`
	Status        string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	AuditRemark   string     `gorm:"type:text" json:"audit_remark"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TrashedAt     *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseApplication) TableName() string {
	return "warehouse_applications"
}

func (x *WarehouseApplication) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
