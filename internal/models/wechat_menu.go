package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WechatMenu represents a WeChat official account custom menu item
type WechatMenu struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ParentID    *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	Type        string     `gorm:"type:varchar(50);default:'click'" json:"type"` // click, view, miniprogram, scancode_push, media_id
	Key         string     `gorm:"type:varchar(255)" json:"key"`
	URL         string     `gorm:"type:varchar(500)" json:"url"`
	AppID       string     `gorm:"type:varchar(255)" json:"appid"`
	PagePath    string     `gorm:"type:varchar(255)" json:"pagepath"`
	BackupURL   string     `gorm:"type:varchar(500)" json:"backup_url"`
	MediaID     string     `gorm:"type:varchar(255)" json:"media_id"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WechatMenu) TableName() string {
	return "wechat_menus"
}

func (x *WechatMenu) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
