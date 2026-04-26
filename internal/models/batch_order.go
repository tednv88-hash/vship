package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BatchOrder links consolidation orders to shipping batches
type BatchOrder struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID  *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	BatchID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"batch_id"`
	OrderID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	CreatedAt time.Time  `json:"created_at"`
}

func (b *BatchOrder) TableName() string {
	return "batch_orders"
}

// BeforeCreate sets UUID before creating
func (b *BatchOrder) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
