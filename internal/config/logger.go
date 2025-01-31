package config

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Define log directory and file path
	logDir := "/var/log"
	if os.Getenv("APP_ENV") != "PROD" {
		logDir = "logs"
	}

	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.Error("Failed to create log directory:", err)
		return logger
	}

	// Open log file
	logPath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("Failed to open log file:", err)
		return logger
	}

	// Set output to both file and stdout
	mw := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(mw)

	// Always use JSON formatter for consistency
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyFunc:  "caller",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", filepath.Base(f.File) + ":" + string(rune(f.Line))
		},
	})

	// Enable caller reporting
	logger.SetReportCaller(true)

	if os.Getenv("APP_ENV") == "PROD" {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
