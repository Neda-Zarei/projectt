package app

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan"
	planD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user"
	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/pkg/database"
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
	db, err := initDB(cfg.DB, log)
	if err != nil {
		return nil, err
	}
	return &app{
		cfg: cfg,
		log: log,
		db:  db,
	}, nil
}

func initDB(c config.DBConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := database.PostgresDSN(
		c.Host, c.Port, c.DBName, c.Schema, c.User, c.Password, c.AppName,
	)
	db, err := database.NewPostgresConnectionWithLogger(dsn, log)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&userD.User{},
		&planD.Plan{},
		&planD.UserPlan{},
		&planD.PlanHistory{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *app) Config() config.Config { return a.cfg }

func (a *app) Logger() *zap.Logger { return a.log }

func (a *app) DB() *gorm.DB { return a.db }

// TODO: replace nil repo with implemented repo
func (a *app) UserService() userP.Service { return user.New(userP.Repo(nil)) }

// TODO: replace nil repo with implemented repo
func (a *app) PlanService() planP.Service { return plan.New(planP.Repo(nil)) }
