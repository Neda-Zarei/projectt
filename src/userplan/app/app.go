package app

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/adapter/repository"
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
	cfg         config.Config
	log         *zap.Logger
	db          *gorm.DB
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	planRepo := repository.NewPlanRepository(db)
	userPlanRepo := repository.NewUserPlanRepository(db)
	priceRepo := repository.NewPriceRepository(db)
	limitationRepo := repository.NewLimitationRepository(db)

	// Initialize services
	userService := user.New(userRepo)
	planService := plan.New(planRepo, userPlanRepo, priceRepo, limitationRepo)

	return &app{
		cfg:         cfg,
		log:         log,
		db:          db,
		userService: userService,
		planService: planService,
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

	//auto-migrate all models
	err = db.AutoMigrate(
		&userD.User{},
		&planD.Plan{},
		&planD.Price{},
		&planD.Limitation{},
		&planD.PlanLimitation{},
		&planD.UserPlan{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *app) Config() config.Config { return a.cfg }

func (a *app) Logger() *zap.Logger { return a.log }

func (a *app) DB() *gorm.DB { return a.db }

func (a *app) UserService() userP.Service { return a.userService }

func (a *app) PlanService() planP.Service { return a.planService }
