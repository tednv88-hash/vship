package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Language represents a supported language
type Language struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`                   // Chinese display name
	EnName      string     `gorm:"type:varchar(100)" json:"enname"`                 // English code name
	LangTo      string     `gorm:"type:varchar(50)" json:"langto"`                  // AI translation code
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"` // active=启用, inactive=禁用
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Language) TableName() string {
	return "languages"
}

func (x *Language) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
