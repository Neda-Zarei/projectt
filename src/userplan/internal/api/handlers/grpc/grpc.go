package grpc

import (
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/pb"
)

var (
	ErrServerCredential = errors.New("failed to generate grpc server credentials")
)

type Server interface {
	Start() error
	Shutdown()
}

type server struct {
	server *grpc.Server
	app    app.App
}

func NewServer(app app.App) (Server, error) {
	cfg := app.Config().GRPC
	var opts []grpc.ServerOption
	if cfg.TLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServerCredential, err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	return &server{
		server: grpc.NewServer(opts...),
		app:    app,
	}, nil
}

func (s *server) Start() error {
	cfg := s.app.Config().GRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	defer func() {
		if err := lis.Close(); err != nil {
			s.app.Logger().Error(err.Error())
		}
	}()
	pb.RegisterUserServiceServer(s.server, newUserServer(s.app.UserService()))
	pb.RegisterPlanServiceServer(s.server, newPlanServer(s.app.PlanService()))
	return s.server.Serve(lis)
}

func (s *server) Shutdown() {
	s.server.GracefulStop()
}
