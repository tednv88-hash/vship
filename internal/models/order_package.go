package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderPackage links packages to consolidation orders
type OrderPackage struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID  *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	PackageID uuid.UUID  `gorm:"type:uuid;not null;index" json:"package_id"`
	CreatedAt time.Time  `json:"created_at"`
}

func (o *OrderPackage) TableName() string {
	return "order_packages"
}

// BeforeCreate sets UUID before creating
func (o *OrderPackage) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
