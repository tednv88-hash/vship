package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscribeMessage represents a subscribe message template
type SubscribeMessage struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	TemplateID  string     `gorm:"type:varchar(255)" json:"template_id"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Content     JSONB      `gorm:"type:jsonb" json:"content"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *SubscribeMessage) TableName() string {
	return "subscribe_messages"
}

func (x *SubscribeMessage) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
