package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func ConnectDB() (*gorm.DB, error) {
	conf := config.Get().Postgres
	db, err := gorm.Open(postgres.Open(prepareDSN()), &gorm.Config{
		Logger: config.NewGormLogger(),
	})

	if err != nil {
		return nil, errors.InternalError("Failed to connect to database", err)
	}

	// Add OpenTelemetry tracing plugin
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, errors.InternalError("Failed to add tracing plugin", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.InternalError("Failed to get database instance", err)
	}

	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeMinutes) * time.Minute)

	return db, nil
}

func CloseDatabaseConnection(db *gorm.DB, logger *config.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.InternalError("Failed to get database instance", err)
	}

	conf := config.Get().Postgres
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.CloseTimeoutSeconds)*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err == nil {
		sqlDB.SetConnMaxLifetime(0)
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
	}

	if err := sqlDB.Close(); err != nil {
		return errors.InternalError("Failed to close database connection", err)
	}

	logger.Info("Database connection closed")

	return nil
}

func IsReady(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

func prepareDSN() string {
	conf := config.Get().Postgres
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		conf.Host, conf.User, conf.Password, conf.Name, conf.Port, conf.TimeZone)
}
