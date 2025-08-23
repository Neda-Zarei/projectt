package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/common"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type userPlanRepository struct {
	db *gorm.DB
}

func NewUserPlanRepository(db *gorm.DB) planP.UserPlanRepository {
	return &userPlanRepository{db: db}
}

func (r *userPlanRepository) AssignPlan(ctx context.Context, userPlan *domain.UserPlan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.UserPlan{}).
			Where("user_id = ? AND status = ?", userPlan.UserID, domain.PlanStatusActive).
			Updates(map[string]interface{}{
				"status":     domain.PlanStatusCanceled,
				"expires_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		if err := tx.Create(userPlan).Error; err != nil {
			return err
		}

		// record history
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

		//record history before update
		history := &domain.PlanHistory{
			UserPlanID: userPlan.ID,
			Action:     domain.PlanActionRenew,
			OldPlanID:  &userPlan.PlanID,
			ChangedAt:  time.Now(),
		}

		//update the plan with new expiration date
		userPlan.ExpiresAt = domain.CalculateExpirationDate(newEndDate, 0) // set to 00:00 of the specified date

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

		userPlan.Status = domain.PlanStatusCanceled
		userPlan.ExpiresAt = time.Now()

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

func (r *userPlanRepository) ExpirePlans(ctx context.Context) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var expiredPlans []domain.UserPlan
		if err := tx.Where("status = ? AND expires_at <= ?",
			domain.PlanStatusActive, time.Now()).Find(&expiredPlans).Error; err != nil {
			return err
		}

		for _, plan := range expiredPlans {
			if err := tx.Model(&plan).Update("status", domain.PlanStatusExpired).Error; err != nil {
				return err
			}

			history := &domain.PlanHistory{
				UserPlanID: plan.ID,
				Action:     "expired",
				OldPlanID:  &plan.PlanID,
				ChangedAt:  time.Now(),
				Metadata: common.JSON{
					"expired_at": time.Now(),
				},
			}
			if err := tx.Create(history).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userPlanRepository) GetExpiringPlans(ctx context.Context, daysThreshold int) ([]*domain.UserPlan, error) {
	var plans []*domain.UserPlan
	thresholdDate := time.Now().AddDate(0, 0, daysThreshold)

	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at <= ? AND expires_at > ?",
			domain.PlanStatusActive, thresholdDate, time.Now()).
		Find(&plans).Error

	return plans, err
}
