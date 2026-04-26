package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Email        string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_email" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Phone        string     `gorm:"type:varchar(50)" json:"phone,omitempty"`
	UserRole     string     `gorm:"type:varchar(50);default:'user'" json:"user_role"`
	Status       string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	ProfilePic   string     `gorm:"type:varchar(500)" json:"profile_pic,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LoggedOutAt  *time.Time `gorm:"type:timestamp with time zone" json:"logged_out_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TrashedAt    *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields  JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
	Tenant       Tenant     `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

// BeforeCreate sets UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.ExtraFields == nil {
		u.ExtraFields = make(JSONB)
	}
	return nil
}
