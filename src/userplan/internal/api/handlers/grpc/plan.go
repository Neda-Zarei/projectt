package grpc

import (
	"context"

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

// TODO: implement pb.UnimplementedPlanServiceServer
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
		DurationDays: int64(res.EndDate.Sub(res.StartDate)),
	}
	return resPB, nil
}

func (s *planServiceServer) RenewUserPlan(ctx context.Context, req *pb.RenewPlanRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
