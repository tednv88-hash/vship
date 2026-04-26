package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailLog represents a log entry for an email message
type EmailLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Email       string     `gorm:"type:varchar(255)" json:"email"`
	Subject     string     `gorm:"type:varchar(500)" json:"subject"`
	Content     string     `gorm:"type:text" json:"content"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *EmailLog) TableName() string {
	return "email_logs"
}

func (x *EmailLog) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
