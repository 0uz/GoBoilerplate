package config

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	envAppEnvironment = "APP_ENV"
	envProdValue      = "PROD"

	defaultSlowThreshold = time.Second // 1 second
	defaultSourceField   = "file"
)

type Logger struct {
	*logrus.Logger
}

var logger *Logger

func NewLogger() *Logger {
	if logger != nil {
		return logger
	}

	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyFunc:  "caller",
		},
	})

	if os.Getenv(envAppEnvironment) == envProdValue {
		l.SetLevel(logrus.InfoLevel)
	} else {
		l.SetLevel(logrus.DebugLevel)
	}

	logger = &Logger{l}
	return logger
}

type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewGormLogger() gormlogger.Interface {
	return &GormLogger{
		SlowThreshold:         defaultSlowThreshold,
		SourceField:           defaultSourceField,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	logger.WithContext(ctx).Infof(s, args...)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logger.WithContext(ctx).Warnf(s, args...)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	logger.WithContext(ctx).Errorf(s, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if l.Debug {
		logger.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}
