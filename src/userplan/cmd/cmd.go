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
	go func() {
		log.Info("server started")
		err = server.Start()
	}()
	
	<-stop
	server.Shutdown()
	if err != nil {
		log.Error("server failed")
	} else {
		log.Info("server shutdown")
	}
	return
}
