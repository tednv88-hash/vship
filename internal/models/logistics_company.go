package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogisticsCompany represents a logistics company (物流公司)
type LogisticsCompany struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
	Code          string     `gorm:"type:varchar(100)" json:"code,omitempty"`
	ContactPerson string     `gorm:"type:varchar(255)" json:"contact_person,omitempty"`
	Phone         string     `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Email         string     `gorm:"type:varchar(255)" json:"email,omitempty"`
	Website       string     `gorm:"type:varchar(500)" json:"website,omitempty"`
	TrackingURL   string     `gorm:"type:varchar(500)" json:"tracking_url,omitempty"`
	Status        string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	SortOrder     int        `gorm:"default:0" json:"sort_order"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TrashedAt     *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (l *LogisticsCompany) TableName() string {
	return "logistics_companies"
}

// BeforeCreate sets UUID before creating
func (l *LogisticsCompany) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.ExtraFields == nil {
		l.ExtraFields = make(JSONB)
	}
	return nil
}
