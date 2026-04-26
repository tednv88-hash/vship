package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderService represents a value-added service applied to an order
type OrderService struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	ServiceID   *uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	ServiceName string     `gorm:"type:varchar(255);not null" json:"service_name"`
	Quantity    int        `gorm:"default:1" json:"quantity"`
	UnitPrice   float64    `gorm:"type:numeric(10,2);default:0" json:"unit_price"`
	TotalPrice  float64    `gorm:"type:numeric(10,2);default:0" json:"total_price"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (o *OrderService) TableName() string {
	return "order_services"
}

// BeforeCreate sets UUID before creating
func (o *OrderService) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
