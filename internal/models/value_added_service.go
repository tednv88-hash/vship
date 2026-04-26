package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ValueAddedService represents a value-added service (增值服務)
type ValueAddedService struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Price       float64    `gorm:"type:numeric(10,2);default:0" json:"price"`
	PriceUnit   string     `gorm:"type:varchar(50);default:'per_item'" json:"price_unit"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (v *ValueAddedService) TableName() string {
	return "value_added_services"
}

// BeforeCreate sets UUID before creating
func (v *ValueAddedService) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.ExtraFields == nil {
		v.ExtraFields = make(JSONB)
	}
	return nil
}
