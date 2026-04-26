package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BatchTracking represents a tracking log entry for a shipping batch
type BatchTracking struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	BatchID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"batch_id"`
	Status      string     `gorm:"type:varchar(50);not null" json:"status"`
	Location    string     `gorm:"type:varchar(255)" json:"location,omitempty"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	OccurredAt  *time.Time `gorm:"type:timestamp with time zone" json:"occurred_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (b *BatchTracking) TableName() string {
	return "batch_tracking_logs"
}

// BeforeCreate sets UUID before creating
func (b *BatchTracking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
