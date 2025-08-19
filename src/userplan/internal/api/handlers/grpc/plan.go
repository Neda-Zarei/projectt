package grpc

import (
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
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
