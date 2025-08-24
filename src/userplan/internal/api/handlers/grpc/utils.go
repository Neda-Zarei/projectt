package grpc

import (
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
)

func UserDomain2Proto(u *userD.User) *pb.User {
	return &pb.User{
		Id:     uint64(u.ID),
		Name:   u.Name,
		Email:  u.Email,
		Phone:  u.Phone,
		Active: u.Active,
	}
}

func UserProto2Domain(u *pb.User) *userD.User {
	return &userD.User{
		Basic: userD.Basic{
			ID: uint(u.Id),
		},
		Email:  u.Email,
		Name:   u.Name,
		Phone:  u.Phone,
		Active: u.Active,
	}
}
