package app

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
)

var (
	ErrNilLogger = errors.New("nil logger")
	ErrInitDB    = errors.New("initial db failed")
)

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
}

func New(cfg config.Config, log *zap.Logger) (App, error) {
	if log == nil {
		return nil, ErrNilLogger
	}
	db, err := initDB(cfg.DB)
	if err != nil {
		return nil, err
	}
	return &app{
		cfg: cfg,
		log: log,
		db:  db,
	}, nil
}

// TODO
func initDB(_ config.DBConfig) (*gorm.DB, error) { return nil, nil }

func (a *app) Config() config.Config { return a.cfg }

func (a *app) Logger() *zap.Logger { return a.log }

func (a *app) DB() *gorm.DB { return a.db }

// TODO: replace nil repo with implemented repo
func (a *app) UserService() userP.Service { return user.New(userP.Repo(nil)) }

// TODO: replace nil repo with implemented repo
func (a *app) PlanService() planP.Service { return plan.New(planP.Repo(nil)) }
