package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/ouz/goauthboilerplate/internal/adapters/api"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	redisCache "github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goauthboilerplate/pkg/errors"

	repoAuth "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/auth"
	repoUser "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/user"

	"github.com/ouz/goauthboilerplate/internal/application/auth"
	"github.com/ouz/goauthboilerplate/internal/application/user"
	"github.com/ouz/goauthboilerplate/internal/config"
	"gorm.io/gorm"
)

var logger *logrus.Logger

func main() {
	logger = config.NewLogger()
	if err := run(); err != nil {
		logger.Error("Application failed to start", "error", err)
	}
}

func run() error {
	if err := config.Load(logger); err != nil {
		return err
	}


	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}
	defer postgres.CloseDatabaseConnection(db, logger)

	// Connect redis cache
	redisClient, err := redisCache.ConnectRedis()
	if err != nil {
		return err
	}
	defer redisCache.CloseRedisClient(redisClient)

	mainRouter := http.NewServeMux()

	mainRouter.Handle("/metrics", promhttp.Handler())

	setupHealthChecks(mainRouter, db)
	setupServiceAndRoutes(mainRouter, db, redisClient)
	mainRouter = addV1Prefix(mainRouter, logger)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Get().App.Port),
		Handler: mainRouter,

		// timeout
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("Server is starting ", "port ", config.Get().App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return errors.GenericError("server shutdown failed: %v", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

func addV1Prefix(r *http.ServeMux, logger *logrus.Logger) *http.ServeMux {
	tokenBucket := middleware.NewTokenBucket(1, 3)
	chain := middleware.Chain(
		middleware.RateLimitMiddleware(tokenBucket),
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)
	v1 := http.NewServeMux()

	v1.Handle("/api/v1/", chain(http.StripPrefix("/api/v1", r)))
	return v1
}

func setupHealthChecks(router *http.ServeMux, db *gorm.DB) {
	router.HandleFunc("/live", livenessHandler)
	router.HandleFunc("/ready", readinessHandler(db))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Live")
}

func readinessHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if postgres.IsReady(db) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Ready")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Not Ready")
		}
	}
}

func setupServiceAndRoutes(mainRouter *http.ServeMux, pgdb *gorm.DB, redisClient *redis.Client) {
	// cache := cache.NewLocalCacheService()
	redisCache := redisCache.NewRedisCacheService(redisClient)
	tx := postgres.NewTransactionManager(pgdb)

	userRepo := repoUser.NewUserRepository(pgdb)
	userService := user.NewUserService(logger, userRepo, redisCache, tx)

	authRepo := repoAuth.NewAuthRepository(pgdb)
	authService := auth.NewAuthService(logger, authRepo, userService, redisCache)

	authHandler := api.NewAuthHandler(logger, authService)
	userHandler := api.NewUserHandler(logger, userService)

	api.SetUpAuthRoutes(mainRouter, authHandler, userHandler, authService)
	api.SetUpUserRoutes(mainRouter, userHandler, authService)

}
