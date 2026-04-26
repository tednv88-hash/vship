package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents a system or business notification (通知)
type Notification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Title       string     `gorm:"type:varchar(500);not null" json:"title"`
	Type        string     `gorm:"type:varchar(50);not null" json:"type"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	TargetType  string     `gorm:"type:varchar(50);default:'all'" json:"target_type"`
	TargetID    *uuid.UUID `gorm:"type:uuid" json:"target_id,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	PublishedAt *time.Time `gorm:"type:timestamp with time zone" json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (n *Notification) TableName() string {
	return "notifications"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.ExtraFields == nil {
		n.ExtraFields = make(JSONB)
	}
	return nil
}
