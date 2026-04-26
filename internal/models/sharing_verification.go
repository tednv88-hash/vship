package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SharingVerification represents a verification record for a sharing order
type SharingVerification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	OrderID     uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	VerifierID  *uuid.UUID `gorm:"type:uuid" json:"verifier_id,omitempty"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *SharingVerification) TableName() string {
	return "sharing_verifications"
}

func (x *SharingVerification) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
