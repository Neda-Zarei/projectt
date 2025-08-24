package plan

import (
	"context"
	"time"

	planD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type service struct {
	planRepo       planP.PlanRepository
	userPlanRepo   planP.UserPlanRepository
	priceRepo      planP.PriceRepository
	limitationRepo planP.LimitationRepository
}

func New(
	planRepo planP.PlanRepository,
	userPlanRepo planP.UserPlanRepository,
	priceRepo planP.PriceRepository,
	limitationRepo planP.LimitationRepository,
) planP.Service {
	return &service{
		planRepo:       planRepo,
		userPlanRepo:   userPlanRepo,
		priceRepo:      priceRepo,
		limitationRepo: limitationRepo,
	}
}

func (s *service) AssignPlan(ctx context.Context, req *planD.AssignPlanRequest) error {
	//soft delete any existing active plan for the user
	existingPlan, err := s.userPlanRepo.GetActiveByUserID(ctx, req.UserID)
	if err == nil && existingPlan != nil {
		if err := s.userPlanRepo.SoftDelete(ctx, existingPlan.ID); err != nil {
			return err
		}
	}

	userPlan := &planD.UserPlan{
		UserID: req.UserID,
		PlanID: req.PlanID,
		ExTime: time.Now().AddDate(0, 1, 0), //default to 1 month, can be customized
	}

	return s.userPlanRepo.Create(ctx, userPlan)
}

func (s *service) GetUserPlan(ctx context.Context, userID uint) (*planD.UserPlan, error) {
	return s.userPlanRepo.GetActiveByUserID(ctx, userID)
}

func (s *service) RenewUserPlan(ctx context.Context, req *planD.RenewPlanRequest) error {
	userPlan, err := s.userPlanRepo.GetActiveByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	userPlan.ExTime = req.EndDate
	return s.userPlanRepo.Update(ctx, userPlan)
}

func (s *service) CancelUserPlan(ctx context.Context, userID uint) error {
	userPlan, err := s.userPlanRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.userPlanRepo.SoftDelete(ctx, userPlan.ID)
}

func (s *service) GetUserPlanHistory(ctx context.Context, userID uint) ([]*planD.UserPlan, error) {
	return s.userPlanRepo.GetUserHistory(ctx, userID)
}

func (s *service) CreatePlan(ctx context.Context, plan *planD.Plan) error {
	return s.planRepo.Create(ctx, plan)
}

func (s *service) GetPlanByID(ctx context.Context, id uint) (*planD.Plan, error) {
	return s.planRepo.GetByID(ctx, id)
}

func (s *service) GetPlanByTitle(ctx context.Context, title string) (*planD.Plan, error) {
	return s.planRepo.GetByTitle(ctx, title)
}

func (s *service) UpdatePlan(ctx context.Context, plan *planD.Plan) error {
	return s.planRepo.Update(ctx, plan)
}

func (s *service) DeletePlan(ctx context.Context, id uint) error {
	return s.planRepo.Delete(ctx, id)
}

func (s *service) ListPlans(ctx context.Context, includeInactive bool) ([]*planD.Plan, error) {
	return s.planRepo.List(ctx, includeInactive)
}

func (s *service) SetPlanPrice(ctx context.Context, planID uint, months int, price int) error {
	existingPrice, err := s.priceRepo.GetByPlanIDAndMonth(ctx, planID, months)
	if err != nil {
		newPrice := &planD.Price{
			PlanID: planID,
			Month:  months,
			Price:  price,
		}
		return s.priceRepo.Create(ctx, newPrice)
	}

	existingPrice.Price = price
	return s.priceRepo.Update(ctx, existingPrice)
}

func (s *service) GetPlanPrices(ctx context.Context, planID uint) ([]*planD.Price, error) {
	return s.priceRepo.GetByPlanID(ctx, planID)
}

func (s *service) CreateLimitation(ctx context.Context, limitation *planD.Limitation) error {
	return s.limitationRepo.Create(ctx, limitation)
}

func (s *service) AssignLimitationToPlan(ctx context.Context, planID, limitationID uint, value int) error {
	planLimitation := &planD.PlanLimitation{
		PlanID:       planID,
		LimitationID: limitationID,
		Value:        value,
	}
	return s.limitationRepo.AssignToPlan(ctx, planLimitation)
}

func (s *service) GetPlanLimitations(ctx context.Context, planID uint) ([]*planD.PlanLimitation, error) {
	return s.limitationRepo.GetPlanLimitations(ctx, planID)
}

// expiration management
func (s *service) ExpirePlans(ctx context.Context) error {
	return s.userPlanRepo.ExpirePlans(ctx)
}

func (s *service) GetExpiringPlans(ctx context.Context, daysThreshold int) ([]*planD.UserPlan, error) {
	return s.userPlanRepo.GetExpiringPlans(ctx, daysThreshold)
}
