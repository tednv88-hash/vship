package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PageCategory represents category page template settings
type PageCategory struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	CategoryStyle string     `gorm:"type:varchar(50);default:'20'" json:"category_style"` // 10=一级分类(大图), 11=一级分类(小图), 20=二级分类
	ShareTitle    string     `gorm:"type:varchar(255)" json:"share_title"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *PageCategory) TableName() string {
	return "page_categories"
}

func (x *PageCategory) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
