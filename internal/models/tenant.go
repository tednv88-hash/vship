package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	if len(bytes) == 0 {
		*j = make(JSONB)
		return nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(bytes, &m); err == nil {
		*j = JSONB(m)
		return nil
	}

	var anyValue interface{}
	if err := json.Unmarshal(bytes, &anyValue); err == nil {
		*j = JSONB{"_data": anyValue}
		return nil
	}

	return json.Unmarshal(bytes, (*map[string]interface{})(j))
}

// Tenant represents a tenant (company/organization)
type Tenant struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Subdomain   string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"subdomain"`
	Plan        string    `gorm:"type:varchar(50);default:'free'" json:"plan"`
	Status      string    `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExtraFields JSONB     `gorm:"type:jsonb;default:'{}'" json:"extra_fields"`
}

func (t *Tenant) TableName() string {
	return "tenants"
}

// BeforeCreate sets UUID before creating
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.ExtraFields == nil {
		t.ExtraFields = make(JSONB)
	}
	return nil
}
