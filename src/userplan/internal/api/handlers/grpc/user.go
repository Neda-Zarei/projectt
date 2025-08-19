package grpc

import (
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	service userP.Service
}

func newUserServer(s userP.Service) pb.UserServiceServer {
	return &userServiceServer{service: s}
}

// TODO: implement pb.UnimplementedUserServiceServer
