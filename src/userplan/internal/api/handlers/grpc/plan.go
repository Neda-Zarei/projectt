package grpc

import (
	"context"
	"time"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
	planD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type planServiceServer struct {
	pb.UnimplementedPlanServiceServer
	service planP.Service
}

func newPlanServer(s planP.Service) pb.PlanServiceServer {
	return &planServiceServer{service: s}
}

// todo: implement pb.UnimplementedPlanServiceServer
func (s *planServiceServer) AssignPlan(ctx context.Context, req *pb.PlanAssignmentRequest) (*pb.Empty, error) {
	reqD := &planD.AssignPlanRequest{
		UserID: uint(req.UserId),
		PlanID: uint(req.PlanId),
	}
	return &pb.Empty{}, s.service.AssignPlan(ctx, reqD)
}

func (s *planServiceServer) GetUserPlan(ctx context.Context, req *pb.UserPlanRequest) (*pb.Plan, error) {
	res, err := s.service.GetUserPlan(ctx, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	resPB := &pb.Plan{
		Id:           uint64(res.PlanID),
		DurationDays: int64(res.ExpiresAt.Sub(res.StartDate) / (24 * time.Hour)),
	}
	return resPB, nil
}

func (s *planServiceServer) RenewUserPlan(ctx context.Context, req *pb.RenewPlanRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *planServiceServer) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.Plan, error) {
	plan := &planD.Plan{
		Name:        req.Plan.Name,
		Description: req.Plan.Description,
		Price:       req.Plan.Price,
		Duration:    int(req.Plan.DurationDays),
		IsActive:    req.Plan.IsActive,
	}

	err := s.service.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:           uint64(plan.ID),
		Name:         plan.Name,
		Description:  plan.Description,
		Price:        plan.Price,
		DurationDays: int64(plan.Duration),
		IsActive:     plan.IsActive,
	}, nil
}

func (s *planServiceServer) GetPlanByID(ctx context.Context, req *pb.PlanIDRequest) (*pb.Plan, error) {
	plan, err := s.service.GetByID(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:           uint64(plan.ID),
		Name:         plan.Name,
		Description:  plan.Description,
		Price:        plan.Price,
		DurationDays: int64(plan.Duration),
		IsActive:     plan.IsActive,
	}, nil
}

func (s *planServiceServer) GetPlanByName(ctx context.Context, req *pb.PlanNameRequest) (*pb.Plan, error) {
	plan, err := s.service.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:           uint64(plan.ID),
		Name:         plan.Name,
		Description:  plan.Description,
		Price:        plan.Price,
		DurationDays: int64(plan.Duration),
		IsActive:     plan.IsActive,
	}, nil
}

func (s *planServiceServer) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.Plan, error) {
	plan := &planD.Plan{
		Name:        req.Plan.Name,
		Description: req.Plan.Description,
		Price:       req.Plan.Price,
		Duration:    int(req.Plan.DurationDays),
		IsActive:    req.Plan.IsActive,
	}
	plan.ID = uint(req.Plan.Id)

	err := s.service.Update(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:           uint64(plan.ID),
		Name:         plan.Name,
		Description:  plan.Description,
		Price:        plan.Price,
		DurationDays: int64(plan.Duration),
		IsActive:     plan.IsActive,
	}, nil
}

func (s *planServiceServer) DeletePlan(ctx context.Context, req *pb.PlanIDRequest) (*pb.Empty, error) {
	err := s.service.Delete(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *planServiceServer) ListPlans(ctx context.Context, req *pb.ListPlansRequest) (*pb.ListPlansResponse, error) {
	plans, err := s.service.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	var pbPlans []*pb.Plan
	for _, plan := range plans {
		pbPlan := &pb.Plan{
			Id:           uint64(plan.ID),
			Name:         plan.Name,
			Description:  plan.Description,
			Price:        plan.Price,
			DurationDays: int64(plan.Duration),
			IsActive:     plan.IsActive,
		}
		pbPlans = append(pbPlans, pbPlan)
	}

	return &pb.ListPlansResponse{
		Plans: pbPlans,
		Total: int64(len(pbPlans)),
	}, nil
}

func (s *planServiceServer) TogglePlanActive(ctx context.Context, req *pb.PlanIDRequest) (*pb.Empty, error) {
	err := s.service.ToggleActive(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
