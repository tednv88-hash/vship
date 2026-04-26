package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseClerk represents a clerk working at a warehouse
type WarehouseClerk struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID uuid.UUID  `gorm:"type:uuid" json:"warehouse_id"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Phone       string     `gorm:"type:varchar(50)" json:"phone"`
	Role        string     `gorm:"type:varchar(50)" json:"role"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseClerk) TableName() string {
	return "warehouse_clerks"
}

func (x *WarehouseClerk) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
