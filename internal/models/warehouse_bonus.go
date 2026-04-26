package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseBonus represents a bonus payment for a warehouse
type WarehouseBonus struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID uuid.UUID  `gorm:"type:uuid" json:"warehouse_id"`
	Amount      float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Description string     `gorm:"type:text" json:"description"`
	Month       string     `gorm:"type:varchar(20)" json:"month"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseBonus) TableName() string {
	return "warehouse_bonuses"
}

func (x *WarehouseBonus) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
