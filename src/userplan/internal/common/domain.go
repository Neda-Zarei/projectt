package common

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type AuditLog struct {
	BaseModel
	UserID     uint   `gorm:"index" json:"user_id"` // admin who performed the action
	Action     string `gorm:"size:255" json:"action"`
	EntityType string `gorm:"size:255" json:"entity_type"`
	EntityID   uint   `json:"entity_id"`
	RequestID  string `gorm:"size:255" json:"request_id"`
	IPAddress  string `gorm:"size:45" json:"ip_address"` // IPv4 or IPv6
	UserAgent  string `gorm:"size:512" json:"user_agent"`
	Metadata   JSON   `gorm:"type:json" json:"metadata"`
}

// for handling JSON data in gorm
type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	// implement scan logic for JSON type
	return nil
}

func (j JSON) Value() (interface{}, error) {
	// implement value logic for JSON type
	return nil, nil
}
