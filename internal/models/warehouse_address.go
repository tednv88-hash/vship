package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseAddress represents a shipping address associated with a warehouse
type WarehouseAddress struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID uuid.UUID  `gorm:"type:uuid" json:"warehouse_id"`
	Label       string     `gorm:"type:varchar(255)" json:"label"`
	Recipient   string     `gorm:"type:varchar(255)" json:"recipient"`
	Phone       string     `gorm:"type:varchar(50)" json:"phone"`
	CountryID   *uuid.UUID `gorm:"type:uuid" json:"country_id,omitempty"`
	Province    string     `gorm:"type:varchar(255)" json:"province"`
	City        string     `gorm:"type:varchar(255)" json:"city"`
	District    string     `gorm:"type:varchar(255)" json:"district"`
	Address     string     `gorm:"type:text" json:"address"`
	PostalCode  string     `gorm:"type:varchar(20)" json:"postal_code"`
	IsDefault   bool       `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WarehouseAddress) TableName() string {
	return "warehouse_addresses"
}

func (x *WarehouseAddress) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
