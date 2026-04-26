package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Feedback represents a user feedback or complaint
type Feedback struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Content     string     `gorm:"type:text" json:"content"`
	Images      JSONB      `gorm:"type:jsonb" json:"images"`
	ContactInfo string     `gorm:"type:varchar(255)" json:"contact_info"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Reply       string     `gorm:"type:text" json:"reply"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Feedback) TableName() string {
	return "feedbacks"
}

func (x *Feedback) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.Images == nil {
		x.Images = make(JSONB)
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
