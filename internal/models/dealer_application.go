package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DealerApplication represents an application to become a dealer
type DealerApplication struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID      uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	RealName    string     `gorm:"type:varchar(255)" json:"real_name"`
	Phone       string     `gorm:"type:varchar(50)" json:"phone"`
	Reason      string     `gorm:"type:text" json:"reason"`
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	AuditRemark string     `gorm:"type:text" json:"audit_remark"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *DealerApplication) TableName() string {
	return "dealer_applications"
}

func (x *DealerApplication) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
