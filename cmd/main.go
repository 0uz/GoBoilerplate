package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ouz/goauthboilerplate/internal/adapters/api"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	redisCache "github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"

	repoAuth "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/auth"
	repoUser "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/user"

	"github.com/ouz/goauthboilerplate/internal/application/auth"
	"github.com/ouz/goauthboilerplate/internal/application/user"
	"github.com/ouz/goauthboilerplate/internal/config"
	"gorm.io/gorm"
)

var logger *slog.Logger

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed to start: %v\n", err)
	}
}

func run() error {
	if err := config.Load(); err != nil {
		return err
	}

	logger = config.NewLogger()

	// Connect postgres database
	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}
	defer postgres.CloseDatabaseConnection(db)

	// Connect redis cache
	redisClient, err := redisCache.ConnectRedis()
	if err != nil {
		return err
	}
	defer redisCache.CloseRedisClient(redisClient)

	mainRouter := http.NewServeMux()

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
		slog.Info("Server is starting", "port", config.Get().App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	slog.Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}

func addV1Prefix(r *http.ServeMux, logger *slog.Logger) *http.ServeMux {
	chain := middleware.Chain(
		middleware.Logging(logger),
		middleware.Recovery(),
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
	userService := user.NewUserService(userRepo, redisCache, tx)

	authRepo := repoAuth.NewAuthRepository(pgdb)
	authService := auth.NewAuthService(authRepo, userService, redisCache)

	authHandler := api.NewAuthHandler(authService)
	userHandler := api.NewUserHandler(userService)

	api.SetUpUserRoutes(mainRouter, authHandler, userHandler, authService)

}
