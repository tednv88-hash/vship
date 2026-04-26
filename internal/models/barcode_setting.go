package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BarcodeSetting represents barcode generation settings
type BarcodeSetting struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	BarcodeType string     `gorm:"type:varchar(50)" json:"barcode_type"`
	Prefix      string     `gorm:"type:varchar(20)" json:"prefix"`
	Length      int        `json:"length"`
	CurrentSeq  int64      `gorm:"default:0" json:"current_seq"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *BarcodeSetting) TableName() string {
	return "barcode_settings"
}

func (x *BarcodeSetting) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
