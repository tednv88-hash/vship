package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Setting represents a system setting key-value pair
type Setting struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID  *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Key       string     `gorm:"type:varchar(255);not null" json:"key"`
	Value     string     `gorm:"type:text" json:"value"`
	Group     string     `gorm:"type:varchar(100)" json:"group,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Setting) TableName() string {
	return "settings"
}

// BeforeCreate sets UUID before creating
func (s *Setting) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
