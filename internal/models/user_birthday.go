package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserBirthday represents a user's birthday information
type UserBirthday struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID        uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	Birthday      *time.Time `gorm:"type:date" json:"birthday,omitempty"`
	LunarBirthday string     `gorm:"type:varchar(20)" json:"lunar_birthday"`
	SendCoupon    bool       `gorm:"default:false" json:"send_coupon"`
	SendPoints    bool       `gorm:"default:false" json:"send_points"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *UserBirthday) TableName() string {
	return "user_birthdays"
}

func (x *UserBirthday) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
