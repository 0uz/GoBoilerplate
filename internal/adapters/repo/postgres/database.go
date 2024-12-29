package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(prepareDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, errors.InternalError("Failed to connect to database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.InternalError("Failed to get database instance", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func CloseDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.InternalError("Failed to get database instance", err)
	}

	// Bağlantıları nazikçe kapat
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err == nil {
		sqlDB.SetConnMaxLifetime(0)
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
	}

	if err := sqlDB.Close(); err != nil {
		return errors.InternalError("Failed to close database connection", err)
	}

	slog.Info("Database connection closed")

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
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
}
