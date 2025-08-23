package database

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	applogger "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/pkg/logger"
)

func NewPostgresConnection(host string, port uint, dbName, schema, user, password, appName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable",
		host, port, user, password, dbName, schema, appName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func NewPostgresConnectionWithLogger(dsn string, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: applogger.NewZapGormLogger(log),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func PostgresDSN(host string, port uint, dbName, schema, user, password, appName string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable",
		host, port, user, password, dbName, schema, appName,
	)
}
