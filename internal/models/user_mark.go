package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserMark represents a mark or tag assigned to a user
type UserMark struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	MarkName    string     `gorm:"type:varchar(255)" json:"mark_name"`
	Remark      string     `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *UserMark) TableName() string {
	return "user_marks"
}

func (x *UserMark) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
