package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
	myhttp "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/handlers/http"
)

func Run(cfg config.Config, log *zap.Logger) error {
	a, err := app.New(cfg, log)
	if err != nil {
		return err
	}

	handler := myhttp.NewHandler(a)
	addr := fmt.Sprintf("%s:%d", a.Config().Server.Host, a.Config().Server.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler.SetupRoutes(),
		ReadTimeout:  time.Duration(a.Config().Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.Config().Server.WriteTimeout) * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		a.Logger().Info("starting HTTP server", zap.String("address", addr))
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-stop:
		a.Logger().Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Error("server failed")
			return err
		}
	}

	shutdownTimeout := time.Duration(a.Config().Server.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed")
		return err
	}

	log.Info("server shutdown complete")
	return nil
}
