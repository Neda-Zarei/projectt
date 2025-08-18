package cmd

import (
	"errors"

	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
)

func Run(cfg config.Config, log *zap.Logger) (err error) {
	a, err := app.New(cfg, log)
	if err != nil {
		return err
	}
	a.Logger().Info("application successfully created")
	a.Logger().Debug("development environment")
	return errors.New("TODO: implement grpc server")
}
