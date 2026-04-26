package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageStatusLog records status transitions for a package
type PackageStatusLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID   *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	PackageID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"package_id"`
	FromStatus string     `gorm:"type:varchar(50);not null" json:"from_status"`
	ToStatus   string     `gorm:"type:varchar(50);not null" json:"to_status"`
	Remark     string     `gorm:"type:text" json:"remark,omitempty"`
	OperatorID uuid.UUID  `gorm:"type:uuid;not null" json:"operator_id"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (psl *PackageStatusLog) TableName() string {
	return "package_status_logs"
}

// BeforeCreate sets UUID before creating
func (psl *PackageStatusLog) BeforeCreate(tx *gorm.DB) error {
	if psl.ID == uuid.Nil {
		psl.ID = uuid.New()
	}
	return nil
}
