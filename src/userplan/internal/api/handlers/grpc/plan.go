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

func (s *planServiceServer) AssignPlan(ctx context.Context, req *pb.PlanAssignmentRequest) (*pb.Empty, error) {
	reqD := &planD.AssignPlanRequest{
		UserID: uint(req.UserId),
		PlanID: uint(req.PlanId),
	}
	return &pb.Empty{}, s.service.AssignPlan(ctx, reqD)
}

func (s *planServiceServer) GetUserPlan(ctx context.Context, req *pb.UserPlanRequest) (*pb.Plan, error) {
	userPlan, err := s.service.GetUserPlan(ctx, uint(req.UserId))
	if err != nil {
		return nil, err
	}

	plan := userPlan.Plan
	resPB := &pb.Plan{
		Id:          uint64(plan.ID),
		Name:        plan.Title,
		Description: "", // add description field to Plan model if needed
		Price:       0,  // get price from prices array based on subscription period
		IsActive:    !planD.IsExpired(userPlan.ExTime),
	}

	//duration in days from ExTime
	if !userPlan.ExTime.IsZero() {
		duration := userPlan.ExTime.Sub(time.Now())
		resPB.DurationDays = int64(duration.Hours() / 24)
	}

	return resPB, nil
}

func (s *planServiceServer) RenewUserPlan(ctx context.Context, req *pb.RenewPlanRequest) (*pb.Empty, error) {
	renewReq := &planD.RenewPlanRequest{
		UserID:  uint(req.UserId),
		EndDate: time.Unix(req.EndDate, 0),
	}
	return &pb.Empty{}, s.service.RenewUserPlan(ctx, renewReq)
}

func (s *planServiceServer) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.Plan, error) {
	plan := &planD.Plan{
		Title:  req.Plan.Name,
		Custom: true,
		PAYG:   false,
	}

	err := s.service.CreatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:          uint64(plan.ID),
		Name:        plan.Title,
		Description: "",
		Price:       0,
		IsActive:    true,
	}, nil
}

func (s *planServiceServer) GetPlanByID(ctx context.Context, req *pb.PlanIDRequest) (*pb.Plan, error) {
	plan, err := s.service.GetPlanByID(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	pbPlan := &pb.Plan{
		Id:       uint64(plan.ID),
		Name:     plan.Title,
		IsActive: plan.Custom || plan.PAYG,
	}

	prices, err := s.service.GetPlanPrices(ctx, plan.ID)
	if err == nil && len(prices) > 0 {
		// use the first price as default, or implement logic to select appropriate price
		pbPlan.Price = float64(prices[0].Price) / 100 // convert from cents to dollars
	}

	return pbPlan, nil
}

func (s *planServiceServer) GetPlanByName(ctx context.Context, req *pb.PlanNameRequest) (*pb.Plan, error) {
	plan, err := s.service.GetPlanByTitle(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	pbPlan := &pb.Plan{
		Id:       uint64(plan.ID),
		Name:     plan.Title,
		IsActive: plan.Custom || plan.PAYG,
	}

	prices, err := s.service.GetPlanPrices(ctx, plan.ID)
	if err == nil && len(prices) > 0 {
		pbPlan.Price = float64(prices[0].Price) / 100
	}

	return pbPlan, nil
}

func (s *planServiceServer) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.Plan, error) {
	plan := &planD.Plan{
		Title:  req.Plan.Name,
		Custom: req.Plan.IsActive,
		PAYG:   false,
	}
	plan.ID = uint(req.Plan.Id)

	err := s.service.UpdatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &pb.Plan{
		Id:       uint64(plan.ID),
		Name:     plan.Title,
		IsActive: plan.Custom || plan.PAYG,
	}, nil
}

func (s *planServiceServer) DeletePlan(ctx context.Context, req *pb.PlanIDRequest) (*pb.Empty, error) {
	err := s.service.DeletePlan(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *planServiceServer) ListPlans(ctx context.Context, req *pb.ListPlansRequest) (*pb.ListPlansResponse, error) {
	plans, err := s.service.ListPlans(ctx, false) // Only active plans
	if err != nil {
		return nil, err
	}

	var pbPlans []*pb.Plan
	for _, plan := range plans {
		pbPlan := &pb.Plan{
			Id:       uint64(plan.ID),
			Name:     plan.Title,
			IsActive: plan.Custom || plan.PAYG,
		}

		prices, err := s.service.GetPlanPrices(ctx, plan.ID)
		if err == nil && len(prices) > 0 {
			pbPlan.Price = float64(prices[0].Price) / 100
		}

		pbPlans = append(pbPlans, pbPlan)
	}

	return &pb.ListPlansResponse{
		Plans: pbPlans,
		Total: int64(len(pbPlans)),
	}, nil
}

func (s *planServiceServer) TogglePlanActive(ctx context.Context, req *pb.PlanIDRequest) (*pb.Empty, error) {
	plan, err := s.service.GetPlanByID(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	//toggle the Custom field (which represents active status in our new model)
	plan.Custom = !plan.Custom
	err = s.service.UpdatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
