package cmd

import (
	"context"
	"time"

	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
)

// runs the plan expiration process
func ExpirePlans(cfg config.Config, log *zap.Logger) error {
	log.Info("Starting plan expiration process")

	a, err := app.New(cfg, log)
	if err != nil {
		log.Error("Failed to create app instance", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	planService := a.PlanService()
	if planService == nil {
		log.Error("Plan service not available")
		return err
	}

	err = planService.ExpirePlans(ctx)
	if err != nil {
		log.Error("Failed to expire plans", zap.Error(err))
		return err
	}

	log.Info("Plan expiration process completed successfully")
	return nil
}

// gets plans that are expiring soon
func GetExpiringPlans(cfg config.Config, log *zap.Logger, daysThreshold int) error {
	log.Info("Getting expiring plans", zap.Int("days_threshold", daysThreshold))

	a, err := app.New(cfg, log)
	if err != nil {
		log.Error("Failed to create app instance", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	planService := a.PlanService()
	if planService == nil {
		log.Error("Plan service not available")
		return err
	}

	plans, err := planService.GetExpiringPlans(ctx, daysThreshold)
	if err != nil {
		log.Error("Failed to get expiring plans", zap.Error(err))
		return err
	}

	log.Info("Found expiring plans", zap.Int("count", len(plans)))
	for _, plan := range plans {
		log.Info("Expiring plan",
			zap.Uint("user_id", plan.UserID),
			zap.Uint("plan_id", plan.PlanID),
			zap.Time("expires_at", plan.ExpiresAt),
		)
	}

	return nil
}
