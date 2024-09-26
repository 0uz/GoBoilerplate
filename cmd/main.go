package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ouz/gobackend/api/routes"
	"github.com/ouz/gobackend/config"
	"github.com/ouz/gobackend/database"
	"github.com/ouz/gobackend/middleware"
	"github.com/ouz/gobackend/pkg/auth"
	"github.com/ouz/gobackend/pkg/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	db, err := database.ConnectDB()
	if err != nil {
		return err
	}
	defer database.CloseDatabaseConnection(db)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, userService)

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(cors.New())
	app.Use(middleware.ErrorHandler)

	api := app.Group("/api/v1")
	routes.SetUpUserRoutes(api, userService, authService)

	go func() {
		if err := app.Listen("localhost:8080"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
