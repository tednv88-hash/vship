package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Consumable represents a packing consumable (耗材)
type Consumable struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Price       float64    `gorm:"type:numeric(10,2);default:0" json:"price"`
	Unit        string     `gorm:"type:varchar(50);default:'個'" json:"unit"`
	Stock       int        `gorm:"default:0" json:"stock"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (c *Consumable) TableName() string {
	return "consumables"
}

// BeforeCreate sets UUID before creating
func (c *Consumable) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExtraFields == nil {
		c.ExtraFields = make(JSONB)
	}
	return nil
}
