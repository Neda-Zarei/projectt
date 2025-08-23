package plan

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/pb"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
)

var (
	ErrPlanNotFound      = errors.New("plan not found")
	ErrPlanAlreadyExists = errors.New("plan already exists")
)

type service struct {
	logger     *zap.Logger
	planClient pb.PlanServiceClient
}

func NewService(logger *zap.Logger, cc *grpc.ClientConn) port.Service {
	return &service{
		logger:     logger,
		planClient: pb.NewPlanServiceClient(cc),
	}
}

func (s *service) CreatePlan(ctx context.Context, plan *domain.Plan) error {
	grpcPlan := &pb.Plan{
		Name:         plan.Name,
		Description:  plan.Description,
		DurationDays: int64(plan.Duration),
		Price:        plan.Price,
		IsActive:     plan.IsActive,
	}

	_, err := s.planClient.CreatePlan(ctx, &pb.CreatePlanRequest{Plan: grpcPlan})
	if err != nil {
		s.logger.Error("Failed to create plan via gRPC", zap.Error(err), zap.String("name", plan.Name))
		return err
	}

	s.logger.Info("Successfully created plan via gRPC", zap.String("name", plan.Name))
	return nil
}

func (s *service) GetPlanByID(ctx context.Context, id uint) (*domain.Plan, error) {
	grpcPlan, err := s.planClient.GetPlanByID(ctx, &pb.PlanIDRequest{Id: uint64(id)})
	if err != nil {
		s.logger.Error("Failed to get plan by ID via gRPC", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	plan := &domain.Plan{
		Name:        grpcPlan.Name,
		Description: grpcPlan.Description,
		Price:       grpcPlan.Price,
		Duration:    int(grpcPlan.DurationDays),
		IsActive:    grpcPlan.IsActive,
	}
	plan.ID = uint(grpcPlan.Id)

	s.logger.Info("Successfully retrieved plan by ID via gRPC", zap.Uint("id", id))
	return plan, nil
}

func (s *service) GetPlanByName(ctx context.Context, name string) (*domain.Plan, error) {
	grpcPlan, err := s.planClient.GetPlanByName(ctx, &pb.PlanNameRequest{Name: name})
	if err != nil {
		s.logger.Error("Failed to get plan by name via gRPC", zap.Error(err), zap.String("name", name))
		return nil, err
	}

	plan := &domain.Plan{
		Name:        grpcPlan.Name,
		Description: grpcPlan.Description,
		Price:       grpcPlan.Price,
		Duration:    int(grpcPlan.DurationDays),
		IsActive:    grpcPlan.IsActive,
	}
	plan.ID = uint(grpcPlan.Id)

	s.logger.Info("Successfully retrieved plan by name via gRPC", zap.String("name", name))
	return plan, nil
}

func (s *service) UpdatePlan(ctx context.Context, plan *domain.Plan) error {
	grpcPlan := &pb.Plan{
		Id:           uint64(plan.ID),
		Name:         plan.Name,
		Description:  plan.Description,
		DurationDays: int64(plan.Duration),
		Price:        plan.Price,
		IsActive:     plan.IsActive,
	}

	_, err := s.planClient.UpdatePlan(ctx, &pb.UpdatePlanRequest{Plan: grpcPlan})
	if err != nil {
		s.logger.Error("Failed to update plan via gRPC", zap.Error(err), zap.Uint("id", plan.ID), zap.String("name", plan.Name))
		return err
	}

	s.logger.Info("Successfully updated plan via gRPC", zap.Uint("id", plan.ID), zap.String("name", plan.Name))
	return nil
}

func (s *service) DeletePlan(ctx context.Context, id uint) error {
	_, err := s.planClient.DeletePlan(ctx, &pb.PlanIDRequest{Id: uint64(id)})
	if err != nil {
		s.logger.Error("Failed to delete plan via gRPC", zap.Error(err), zap.Uint("id", id))
		return err
	}

	s.logger.Info("Successfully deleted plan via gRPC", zap.Uint("id", id))
	return nil
}

func (s *service) ListPlans(ctx context.Context, limit, offset int) ([]*domain.Plan, error) {
	response, err := s.planClient.ListPlans(ctx, &pb.ListPlansRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		s.logger.Error("Failed to list plans via gRPC", zap.Error(err))
		return nil, err
	}

	var plans []*domain.Plan
	for _, grpcPlan := range response.Plans {
		plan := &domain.Plan{
			Name:        grpcPlan.Name,
			Description: grpcPlan.Description,
			Price:       grpcPlan.Price,
			Duration:    int(grpcPlan.DurationDays),
			IsActive:    grpcPlan.IsActive,
		}
		plan.ID = uint(grpcPlan.Id)
		plans = append(plans, plan)
	}

	s.logger.Info("Successfully listed plans via gRPC", zap.Int("count", len(plans)))
	return plans, nil
}

func (s *service) TogglePlanActive(ctx context.Context, id uint) error {
	_, err := s.planClient.TogglePlanActive(ctx, &pb.PlanIDRequest{Id: uint64(id)})
	if err != nil {
		s.logger.Error("Failed to toggle plan active status via gRPC", zap.Error(err), zap.Uint("id", id))
		return err
	}

	s.logger.Info("Successfully toggled plan active status via gRPC", zap.Uint("id", id))
	return nil
}
