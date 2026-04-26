package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlindBoxDraw represents a single draw in a blind box activity
type BlindBoxDraw struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ActivityID  uuid.UUID  `gorm:"type:uuid" json:"activity_id"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	PrizeInfo   JSONB      `gorm:"type:jsonb" json:"prize_info"`
	CostPoints  int        `json:"cost_points"`
	CreatedAt   time.Time  `json:"created_at"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *BlindBoxDraw) TableName() string {
	return "blind_box_draws"
}

func (x *BlindBoxDraw) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
