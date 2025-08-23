package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/api/handlers/grpc"
)

func Run(cfg config.Config, log *zap.Logger) (err error) {
	a, err := app.New(cfg, log)
	if err != nil {
		return err
	}
	server, err := grpc.NewServer(a)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	serverErr := make(chan error, 1)
	go func() {
		log.Info("server started")
		serverErr <- server.Start()
	}()

	select {
	case <-stop:
		a.Logger().Info("shutdown signal received")
	case err := <-serverErr:
		a.Logger().Error("server failed")
		return err
	}

	server.Shutdown()
	log.Info("server shutdown complete")
	return nil
}
