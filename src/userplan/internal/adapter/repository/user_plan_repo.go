// designed for use by grpc handlers in the userplan service
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
)

type UserPlanRepository interface {
	AssignPlan(ctx context.Context, userPlan *domain.UserPlan) error
	GetByUserID(ctx context.Context, userID uint) (*domain.UserPlan, error)
	RenewPlan(ctx context.Context, userID uint, newEndDate time.Time) error
	CancelPlan(ctx context.Context, userID uint) error
	GetHistory(ctx context.Context, userID uint) ([]*domain.PlanHistory, error)
	RecordHistory(ctx context.Context, history *domain.PlanHistory) error
}

type userPlanRepository struct {
	db *gorm.DB
}

func NewUserPlanRepository(db *gorm.DB) UserPlanRepository {
	return &userPlanRepository{db: db}
}

func (r *userPlanRepository) AssignPlan(ctx context.Context, userPlan *domain.UserPlan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//canceling any existing active plan
		if err := tx.Model(&domain.UserPlan{}).
			Where("user_id = ? AND status = ?", userPlan.UserID, domain.PlanStatusActive).
			Updates(map[string]interface{}{
				"status":   domain.PlanStatusCanceled,
				"end_date": time.Now(),
			}).Error; err != nil {
			return err
		}

		//a new plan assignment
		if err := tx.Create(userPlan).Error; err != nil {
			return err
		}

		//recording history
		history := &domain.PlanHistory{
			UserPlanID: userPlan.ID,
			Action:     domain.PlanActionAssign,
			NewPlanID:  &userPlan.PlanID,
			ChangedAt:  time.Now(),
		}
		return tx.Create(history).Error
	})
}

func (r *userPlanRepository) GetByUserID(ctx context.Context, userID uint) (*domain.UserPlan, error) {
	var userPlan domain.UserPlan
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, domain.PlanStatusActive).
		First(&userPlan).Error
	return &userPlan, err
}

func (r *userPlanRepository) RenewPlan(ctx context.Context, userID uint, newEndDate time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userPlan domain.UserPlan
		if err := tx.Where("user_id = ? AND status = ?", userID, domain.PlanStatusActive).
			First(&userPlan).Error; err != nil {
			return err
		}

		//recoridng history before update
		history := &domain.PlanHistory{
			UserPlanID: userPlan.ID,
			Action:     domain.PlanActionRenew,
			OldPlanID:  &userPlan.PlanID,
			ChangedAt:  time.Now(),
		}

		//updating the plan
		userPlan.EndDate = newEndDate
		userPlan.LastRenewalAt = time.Now()
		userPlan.NextRenewalAt = newEndDate.Add(-7 * 24 * time.Hour) // 7 days before end ??

		if err := tx.Save(&userPlan).Error; err != nil {
			return err
		}

		return tx.Create(history).Error
	})
}

func (r *userPlanRepository) CancelPlan(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userPlan domain.UserPlan
		if err := tx.Where("user_id = ? AND status = ?", userID, domain.PlanStatusActive).
			First(&userPlan).Error; err != nil {
			return err
		}

		history := &domain.PlanHistory{
			UserPlanID: userPlan.ID,
			Action:     domain.PlanActionCancel,
			OldPlanID:  &userPlan.PlanID,
			ChangedAt:  time.Now(),
		}

		//updating plan status
		userPlan.Status = domain.PlanStatusCanceled
		userPlan.EndDate = time.Now()

		if err := tx.Save(&userPlan).Error; err != nil {
			return err
		}

		return tx.Create(history).Error
	})
}

func (r *userPlanRepository) GetHistory(ctx context.Context, userID uint) ([]*domain.PlanHistory, error) {
	var history []*domain.PlanHistory
	err := r.db.WithContext(ctx).
		Joins("JOIN user_plans ON user_plans.id = plan_histories.user_plan_id").
		Where("user_plans.user_id = ?", userID).
		Order("plan_histories.changed_at DESC").
		Find(&history).Error
	return history, err
}

func (r *userPlanRepository) RecordHistory(ctx context.Context, history *domain.PlanHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}
