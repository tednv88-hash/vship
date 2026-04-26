package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebMenu represents a website navigation menu item
type WebMenu struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	ParentID    *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	Type        int        `gorm:"default:10" json:"type"` // 10=单页, 20=列表, 30=关于我们, 40=仓库地址
	LinkID      string     `gorm:"type:varchar(255)" json:"link_id"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WebMenu) TableName() string {
	return "web_menus"
}

func (x *WebMenu) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
