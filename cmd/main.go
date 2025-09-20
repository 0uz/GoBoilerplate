package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	errs "errors"

	"github.com/ouz/goauthboilerplate/internal/adapters/api"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	redisCache "github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goauthboilerplate/internal/observability"
	"github.com/ouz/goauthboilerplate/pkg/errors"

	repoAuth "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/auth"
	repoUser "github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres/user"

	"github.com/ouz/goauthboilerplate/internal/application/auth"
	"github.com/ouz/goauthboilerplate/internal/application/user"
	"github.com/ouz/goauthboilerplate/internal/config"
	"gorm.io/gorm"
)

var logger *config.Logger

func main() {
	logger = config.NewLogger()
	if err := run(); err != nil {
		logger.Error("Application failed to start", "error", err)
	}
}

func run() error {
	ctx := context.Background()

	if err := config.Load(logger); err != nil {
		return err
	}

	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}

	defer func() {
		if err := postgres.CloseDatabaseConnection(db, logger); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	redisClient, err := redisCache.ConnectRedis()
	if err != nil {
		return err
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection", "error", err)
		}
	}()

	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		return err
	}

	// Setup OTEL slog bridge
	_ = observability.SetupOTelSlog()

	defer func() {
		err = errs.Join(err, otelShutdown(context.Background()))
	}()

	response.InitResponseLogger(logger)

	businessRouter := http.NewServeMux()
	setupServiceAndRoutes(businessRouter, db, redisClient)

	mainRouter := createFinalRouter(businessRouter, db, logger)

	// Filter function to skip tracing for health check endpoints
	skipPaths := map[string]bool{
		"/":        true,
		"/ready":   true,
		"/health":  true,
		"/ping":    true,
		"/metrics": true,
	}

	filter := func(req *http.Request) bool {
		return !skipPaths[req.URL.Path]
	}

	mainRouterWithOTel := otelhttp.NewHandler(mainRouter, "go-auth-boilerplate",
		otelhttp.WithFilter(filter),
		otelhttp.WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
		otelhttp.WithPublicEndpoint(),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Get().App.Port),
		Handler: mainRouterWithOTel,

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

	logger.Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return errors.GenericError("server shutdown failed: %v", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

func createFinalRouter(businessRouter *http.ServeMux, db *gorm.DB, logger *config.Logger) *http.ServeMux {
	chain := middleware.Chain(
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)

	finalRouter := http.NewServeMux()

	finalRouter.Handle("/metrics", promhttp.Handler())
	finalRouter.HandleFunc("/live", livenessHandler)
	finalRouter.HandleFunc("/ready", readinessHandler(db))

	finalRouter.Handle("/api/v1/", chain(http.StripPrefix("/api/v1", businessRouter)))

	return finalRouter
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
