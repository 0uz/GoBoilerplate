package config

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	if os.Getenv("APP_ENV") == "PROD" {

		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})

		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			PadLevelText:    true,
			FullTimestamp:   true,
			ForceColors:     true,
		})

		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetOutput(os.Stdout)
	return logger
}
