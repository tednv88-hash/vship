package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddressBook represents a member's saved address (地址簿)
type AddressBook struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	MemberID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"member_id"`
	RecipientName string     `gorm:"type:varchar(200);not null" json:"recipient_name"`
	Phone         string     `gorm:"type:varchar(50);not null" json:"phone"`
	CountryID     *uuid.UUID `gorm:"type:uuid;index" json:"country_id,omitempty"`
	CountryName   string     `gorm:"-" json:"country_name,omitempty"`
	Province      string     `gorm:"type:varchar(200)" json:"province,omitempty"`
	City          string     `gorm:"type:varchar(200)" json:"city,omitempty"`
	District      string     `gorm:"type:varchar(200)" json:"district,omitempty"`
	Address       string     `gorm:"type:text;not null" json:"address"`
	PostalCode    string     `gorm:"type:varchar(50)" json:"postal_code,omitempty"`
	IsDefault     bool       `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TrashedAt     *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields   JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (ab *AddressBook) TableName() string {
	return "address_books"
}

// BeforeCreate sets UUID and ExtraFields before creating
func (ab *AddressBook) BeforeCreate(tx *gorm.DB) error {
	if ab.ID == uuid.Nil {
		ab.ID = uuid.New()
	}
	if ab.ExtraFields == nil {
		ab.ExtraFields = make(JSONB)
	}
	return nil
}
