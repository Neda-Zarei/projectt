package app

import (
	"errors"

	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
)

var ErrNilLogger = errors.New("nil logger")

type App interface {
	Config() config.Config
	Logger() *zap.Logger
}

type app struct {
	cfg config.Config
	log *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) (App, error) {
	if log == nil {
		return nil, ErrNilLogger
	}
	return &app{
		cfg: cfg,
		log: log,
	}, nil
}

func (a *app) Config() config.Config { return a.cfg }

func (a *app) Logger() *zap.Logger { return a.log }
