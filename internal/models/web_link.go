package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebLink represents a friendly link (友情链接)
type WebLink struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	URL         string     `gorm:"type:varchar(500)" json:"url"`
	ImageURL    string     `gorm:"type:varchar(500)" json:"image_url"`
	LinkType    string     `gorm:"type:varchar(50);default:'image'" json:"link_type"` // image, text
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	Status      string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *WebLink) TableName() string {
	return "web_links"
}

func (x *WebLink) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
