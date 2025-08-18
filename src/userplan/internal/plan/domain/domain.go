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

type Plan struct {
	common.BaseModel
	Name        string      `gorm:"size:100;uniqueIndex" json:"name"`
	Description string      `gorm:"size:500" json:"description"`
	Price       float64     `gorm:"type:decimal(10,2)" json:"price"`
	Duration    int         `gorm:"default:30" json:"duration"` // in days
	Features    common.JSON `gorm:"type:json" json:"features"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
}

type UserPlan struct {
	common.BaseModel
	UserID         uint      `gorm:"index" json:"user_id"`
	PlanID         uint      `gorm:"index" json:"plan_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsAutoRenew    bool      `gorm:"default:true" json:"is_auto_renew"`
	LastRenewalAt  time.Time `json:"last_renewal_at"`
	NextRenewalAt  time.Time `json:"next_renewal_at"`
	Status         string    `gorm:"size:20;default:active" json:"status"`
	PaymentGateway string    `gorm:"size:50" json:"payment_gateway"`
	TransactionID  string    `gorm:"size:255" json:"transaction_id"`
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
