package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SigninLog represents a daily check-in record
type SigninLog struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_signin_logs_tenant_user_date" json:"tenant_id,omitempty"`
	UserID          uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_signin_logs_tenant_user_date" json:"user_id"`
	SigninDate      string     `gorm:"type:varchar(10);uniqueIndex:idx_signin_logs_tenant_user_date" json:"signin_date"`
	PointsEarned    int        `json:"points_earned"`
	ConsecutiveDays int        `json:"consecutive_days"`
	CreatedAt       time.Time  `json:"created_at"`
	ExtraFields     JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *SigninLog) TableName() string {
	return "signin_logs"
}

func (x *SigninLog) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
