package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SmsLog represents a log entry for an SMS message
type SmsLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Phone       string     `gorm:"type:varchar(50)" json:"phone"`
	Content     string     `gorm:"type:text" json:"content"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *SmsLog) TableName() string {
	return "sms_logs"
}

func (x *SmsLog) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
