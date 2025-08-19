package cmd

import (
	"context"
	"fmt"
	"net/http"
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
	a.Logger().Info("application successfully created")

	if a.Config().DevEnv {
		a.Logger().Debug("development environment")
	}

	handler := myhttp.NewHandler(a)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.Config().Server.Host, a.Config().Server.Port),
		Handler:      handler.SetupRoutes(),
		ReadTimeout:  time.Duration(a.Config().Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.Config().Server.WriteTimeout) * time.Second,
	}

	a.Logger().Info("starting HTTP server",
		zap.String("host", a.Config().Server.Host),
		zap.Int("port", a.Config().Server.Port),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger().Fatal("HTTP server error", zap.Error(err))
		}
	}()

	//wait for shutdown sig
	<-context.Background().Done()

	a.Logger().Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
