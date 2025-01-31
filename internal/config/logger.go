package config

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyFunc:  "caller",
		},
	})

	if os.Getenv("APP_ENV") == "PROD" {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
