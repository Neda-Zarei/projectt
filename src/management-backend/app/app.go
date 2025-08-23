package app

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/adapter/repository"
	user "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/pkg/database"
)

var ErrNilLogger = errors.New("nil logger")

type App interface {
	Config() config.Config
	Logger() *zap.Logger
	DB() *gorm.DB
	UserService() userP.Service
	PlanService() planP.Service
}

type app struct {
	cfg config.Config
	log *zap.Logger
	db  *gorm.DB
	cc  *grpc.ClientConn

	userService userP.Service
	planService planP.Service
}

func New(cfg config.Config, log *zap.Logger) (App, error) {
	if log == nil {
		return nil, ErrNilLogger
	}
	db, err := initDB(cfg.DB, log)
	if err != nil {
		return nil, err
	}
	cc, err := newGRPCClientConn(cfg.UserPlanService)
	if err != nil {
		return nil, err
	}
	return &app{
		cfg: cfg,
		log: log,
		db:  db,
		cc:  cc,
	}, nil
}

func (a *app) Config() config.Config { return a.cfg }

func (a *app) Logger() *zap.Logger { return a.log }

func (a *app) DB() *gorm.DB { return a.db }

func (a *app) UserService() userP.Service {
	if a.userService == nil {
		a.userService = user.NewService(repository.NewUserRepository(a.db), a.cc)
	}
	return a.userService
}

func (a *app) PlanService() planP.Service {
	if a.planService == nil {
		a.planService = plan.NewService(a.log, a.cc)
	}
	return a.planService
}

func initDB(c config.DBConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := database.PostgresDSN(
		c.Host, c.Port, c.DBName, c.Schema, c.User, c.Password, c.AppName,
	)
	db, err := database.NewPostgresConnectionWithLogger(dsn, log)
	if err != nil {
		return nil, err
	}

	// only migrate admin tables
	// plan data will be managed by the userplan service via gRPC
	err = db.AutoMigrate(
	// add your admin domain models here
	// &admin.Domain.AdminUser{},
	)
	if err != nil {
		return nil, err
	}

	// todo: initial admin setup
	return db, nil
}

func newGRPCClientConn(cfg config.UserPlanServiceConfig) (*grpc.ClientConn, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	cc, err := grpc.NewClient(addr, creds)
	if err != nil {
		return nil, err
	}
	return cc, nil
}
