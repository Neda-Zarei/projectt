package domain

import (
	"time"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/common"
)

const (
	PlanStatusActive   = "active"
	PlanStatusExpired  = "expired"
	PlanStatusCanceled = "canceled"
)

const (
	PlanActionAssign    = "assign"
	PlanActionRenew     = "renew"
	PlanActionCancel    = "cancel"
	PlanActionUpgrade   = "upgrade"
	PlanActionDowngrade = "downgrade"
)

// subscription plan template
type Plan struct {
	common.BaseModel
	Name        string      `gorm:"size:100;uniqueIndex" json:"name"`
	Description string      `gorm:"size:500" json:"description"`
	Price       float64     `gorm:"type:decimal(10,2)" json:"price"`
	Duration    int         `gorm:"default:30" json:"duration"` // in days
	Features    common.JSON `gorm:"type:json" json:"features"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
}

// user's active subscriptions
type UserPlan struct {
	common.BaseModel
	UserID    uint      `gorm:"index" json:"user_id"`
	PlanID    uint      `gorm:"index" json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"` // expiration date at 00:00 of the day
	Status    string    `gorm:"size:20;default:active" json:"status"`
}

// tracks changes to user plans
type PlanHistory struct {
	common.BaseModel
	UserPlanID uint        `gorm:"index" json:"user_plan_id"`
	Action     string      `gorm:"size:50" json:"action"`
	OldPlanID  *uint       `json:"old_plan_id,omitempty"` // nullable for new assignments
	NewPlanID  *uint       `json:"new_plan_id,omitempty"` // nullable for cancellations
	ChangedAt  time.Time   `json:"changed_at"`
	Metadata   common.JSON `gorm:"type:json" json:"metadata"`
}

// DTOs
type AssignPlanRequest struct {
	UserID uint
	PlanID uint
}

type RenewPlanRequest struct {
	UserID  uint
	EndDate time.Time
}

func CalculateExpirationDate(startDate time.Time, durationDays int) time.Time {
	expirationDate := startDate.AddDate(0, 0, durationDays)

	//setting the time to 00:00:00 of the expiration day
	return time.Date(expirationDate.Year(), expirationDate.Month(), expirationDate.Day(), 0, 0, 0, 0, expirationDate.Location())
}

func IsExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

func IsExpiringSoon(expiresAt time.Time, daysThreshold int) bool {
	thresholdDate := time.Now().AddDate(0, 0, daysThreshold)
	return expiresAt.Before(thresholdDate) && !IsExpired(expiresAt)
}
