package database

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/ouz/gobackend/config"
	"github.com/ouz/gobackend/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, errors.InternalError("Failed to connect to database", err)
	}

	return db, nil
}

func CloseDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.InternalError("Failed to get database instance", err)
	}

	if err := sqlDB.Close(); err != nil {
		return errors.InternalError("Failed to close database connection", err)
	}

	log.Info("Database connection closed")

	return nil
}

func IsReady(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
