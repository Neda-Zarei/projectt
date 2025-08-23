package domain

import "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/common"

type Plan struct {
	common.BaseModel
	Name        string  `gorm:"size:100;uniqueIndex" json:"name" validate:"required"`
	Description string  `gorm:"size:500" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price" validate:"gte=0"`
	Duration    int     `gorm:"default:30" json:"duration" validate:"gte=30"` // in days
	IsActive    bool    `gorm:"default:true" json:"is_active"`
}
