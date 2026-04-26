package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File represents an uploaded file
type File struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	GroupID     *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`
	FileName    string     `gorm:"type:varchar(500)" json:"file_name"`
	FilePath    string     `gorm:"type:varchar(1000)" json:"file_path"`
	FileURL     string     `gorm:"type:varchar(1000)" json:"file_url"`
	FileType    string     `gorm:"type:varchar(100)" json:"file_type"`
	FileSize    int64      `json:"file_size"`
	StorageType string     `gorm:"type:varchar(50)" json:"storage_type"`
	UploaderID  *uuid.UUID `gorm:"type:uuid" json:"uploader_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TrashedAt   *time.Time `gorm:"type:timestamp with time zone" json:"trashed_at,omitempty"`
	ExtraFields JSONB      `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (x *File) TableName() string {
	return "files"
}

func (x *File) BeforeCreate(tx *gorm.DB) error {
	if x.ID == uuid.Nil {
		x.ID = uuid.New()
	}
	if x.ExtraFields == nil {
		x.ExtraFields = make(JSONB)
	}
	return nil
}
