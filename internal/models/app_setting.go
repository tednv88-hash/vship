package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppSetting represents application-level settings
type AppSetting struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	SettingType string     `gorm:"type:varchar(50)" json:"setting_type"`
	Config      JSONB      `gorm:"type:jsonb" json:"config"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *AppSetting) TableName() string {
	return "app_settings"
}

func (x *AppSetting) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
