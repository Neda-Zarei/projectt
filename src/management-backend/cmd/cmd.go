package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
	myhttp "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/handlers/http"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/pkg/database"
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

	db, err := database.NewPostgresConnection(
		a.Config().DB.Host,
		a.Config().DB.Port,
		a.Config().DB.DBName,
		a.Config().DB.Schema,
		a.Config().DB.User,
		a.Config().DB.Password,
		a.Config().DB.AppName,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	//auto migrate tables
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
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

	//wait for shutdown signal
	<-context.Background().Done()

	a.Logger().Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func autoMigrate(db *gorm.DB) error {
	// add models here for auto migration
	// return db.AutoMigrate(
	//     &user.domain.AdminUser{},
	//     &plan.domain.Plan{},
	// )
	return nil
}
