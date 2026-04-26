package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarehouseShelf represents a shelf or slot within a warehouse
type WarehouseShelf struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	WarehouseID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"warehouse_id"`
	ShelfCode    string     `gorm:"type:varchar(100);not null" json:"shelf_code"`
	Zone         string     `gorm:"type:varchar(100)" json:"zone,omitempty"`
	Capacity     int        `gorm:"default:0" json:"capacity"`
	UsedCapacity int        `gorm:"default:0" json:"used_capacity"`
	Status       string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ExtraFields  JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (ws *WarehouseShelf) TableName() string {
	return "warehouse_shelves"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (ws *WarehouseShelf) BeforeCreate(tx *gorm.DB) error {
	if ws.ID == uuid.Nil {
		ws.ID = uuid.New()
	}
	if ws.ExtraFields == nil {
		ws.ExtraFields = make(JSONB)
	}
	return nil
}
