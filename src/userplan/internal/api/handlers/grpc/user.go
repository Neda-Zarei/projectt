package grpc

import (
	"context"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/pkg/util"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	service userP.Service
}

func newUserServer(s userP.Service) pb.UserServiceServer {
	return &userServiceServer{service: s}
}

func (s *userServiceServer) ListUsers(ctx context.Context, uf *pb.UserFilter) (*pb.PaginatedUsers, error) {
	res, err := s.service.ListUsers(ctx, &userD.UserFilter{
		Name:  uf.Name,
		Email: uf.Email,
		Phone: uf.Phone,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PaginatedUsers{
		Users: util.Map(res.Users, UserDomain2Proto),
		Total: res.Total,
		Size:  res.Size,
		Page:  res.Page,
	}, nil
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.service.CreateUser(ctx, UserProto2Domain(req.User))
}

func (s *userServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.service.UpdateUser(ctx, UserProto2Domain(req.User))
}

func (s *userServiceServer) SetUserActive(ctx context.Context, ua *pb.UserActivationRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.service.SetUserActive(ctx, &userD.UserActivation{
		UserID: uint(ua.UserId), Active: ua.Active})
}
