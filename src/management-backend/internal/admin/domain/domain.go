package domain

import (
	"time"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/common"
)

type AdminUser struct {
	common.BaseModel
	Email        string    `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash string    `gorm:"size:255" json:"passwordhash"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	LastLogin    time.Time `json:"last_login"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Role         string    `gorm:"size:50;default:admin" json:"role"`
}
