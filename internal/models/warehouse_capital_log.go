package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseCapitalLog represents a capital transaction log for a warehouse
type WarehouseCapitalLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID uuid.UUID  `gorm:"type:uuid" json:"warehouse_id"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Balance     float64    `gorm:"type:numeric(12,2)" json:"balance"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Description string     `gorm:"type:text" json:"description"`
	RelatedID   *uuid.UUID `gorm:"type:uuid" json:"related_id,omitempty"`
	RelatedType string     `gorm:"type:varchar(100)" json:"related_type"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseCapitalLog) TableName() string {
	return "warehouse_capital_logs"
}

func (x *WarehouseCapitalLog) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
