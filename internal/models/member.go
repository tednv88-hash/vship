package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Member represents a registered member (會員)
type Member struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name            string     `gorm:"type:varchar(200);not null" json:"name"`
	Email           string     `gorm:"type:varchar(200)" json:"email,omitempty"`
	Phone           string     `gorm:"type:varchar(50)" json:"phone,omitempty"`
	MemberLevelID   *uuid.UUID `gorm:"type:uuid;index" json:"member_level_id,omitempty"`
	MemberLevelName string     `gorm:"-" json:"member_level_name,omitempty"`
	Balance         float64    `gorm:"type:decimal(12,2);default:0" json:"balance"`
	AvatarURL       string     `gorm:"type:text" json:"avatar_url,omitempty"`
	Status          string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	Remark          string     `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TrashedAt       *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields     JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (m *Member) TableName() string {
	return "members"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.ExtraFields == nil {
		m.ExtraFields = make(JSONB)
	}
	return nil
}
