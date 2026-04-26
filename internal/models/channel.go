package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Channel represents a sales or distribution channel
type Channel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Code        string     `gorm:"type:varchar(100)" json:"code"`
	Type        string     `gorm:"type:varchar(50)" json:"type"`
	Config      JSONB      `gorm:"type:jsonb" json:"config"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *Channel) TableName() string {
	return "channels"
}

func (x *Channel) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
