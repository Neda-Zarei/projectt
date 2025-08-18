package domain

import "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/common"

type Plan struct {
	common.BaseModel
	Name        string  `gorm:"size:100;uniqueIndex" json:"name"`
	Description string  `gorm:"size:500" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`
	Duration    int     `gorm:"default:30" json:"duration"` // in days
	IsActive    bool    `gorm:"default:true" json:"is_active"`
}
