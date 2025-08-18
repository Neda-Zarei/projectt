package domain

import (
	"time"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/common"
)

// end users who can be assigned plans
type User struct {
	common.BaseModel
	ExternalID  string    `gorm:"size:255;uniqueIndex" json:"external_id"`
	Email       string    `gorm:"size:255;index" json:"email"`
	Phone       string    `gorm:"size:20;index" json:"phone"`
	FirstName   string    `gorm:"size:100" json:"first_name"`
	LastName    string    `gorm:"size:100" json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}
