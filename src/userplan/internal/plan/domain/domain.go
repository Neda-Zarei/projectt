package domain

import (
	"time"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/common"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
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

type BasicID struct {
	ID uint `gorm:"primarykey"`
}

type BasicWithSoftDelete struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"<-:create;"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Plan struct {
	BasicID
	Title       string `gorm:"not null"`
	Prices      []Price
	Limitations []Limitation `gorm:"many2many:plan_limitations;"`
	Custom      bool         `gorm:"not null;default:true"`
	PAYG        bool         `gorm:"not null;default:false"`
}

type Price struct {
	PlanID uint `gorm:"primaryKey"`
	Month  int  `gorm:"primaryKey"`
	Price  int  `gorm:"not null;default:1000000"`
}

type Limitation struct {
	BasicID
	Title string `gorm:"not null;unique"`
}

type PlanLimitation struct {
	PlanID       uint `gorm:"primaryKey"`
	LimitationID uint `gorm:"primaryKey"`
	Limitation   Limitation
	Value        int `gorm:"default:1"`
}

type UserPlan struct {
	BasicWithSoftDelete
	PlanID uint `gorm:"primaryKey"`
	Plan   Plan
	UserID uint `gorm:"primaryKey"`
	User   domain.User
	ExTime time.Time
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
