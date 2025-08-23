package grpc

import (
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/common"
	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
)

func UserDomain2Porto(u *userD.User) *pb.User {
	return &pb.User{
		Id:     uint64(u.ID),
		Name:   u.FirstName + " " + u.LastName,
		Email:  u.Email,
		Phone:  u.Phone,
		Active: u.IsActive,
	}
}

func UserProto2Domain(u *pb.User) *userD.User {
	return &userD.User{
		BaseModel: common.BaseModel{
			ID: uint(u.Id),
		},
		ExternalID: "", // ?
		Email:      u.Email,
		Phone:      u.Phone,
		FirstName:  u.Name, // !
		LastName:   u.Name, // !
		IsActive:   u.Active,
	}
}
