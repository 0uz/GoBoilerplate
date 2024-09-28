package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/api/handlers"
	"github.com/ouz/gobackend/middleware"
	"github.com/ouz/gobackend/pkg/auth"
	entity "github.com/ouz/gobackend/pkg/entities"
	"github.com/ouz/gobackend/pkg/user"
)

func SetUpUserRoutes(app fiber.Router, userService user.Service, authService auth.Service) {
	user := app.Group("/user")
	user.Get("/login", handlers.LoginUser(userService, authService))
	user.Post("/register", handlers.RegisterUser(userService, authService))
	user.Post("/anonymous", handlers.RegisterAnonymousUser(userService, authService))
	user.Post("/refresh", handlers.RefreshAccessToken(userService, authService))

	protectedUsers := app.Group("/user")
	protectedUsers.Use(middleware.Protected(authService))
	protectedUsers.Use(middleware.HasRoles(entity.UserRoleUser))
	protectedUsers.Get("/all", handlers.AllUsers(userService))
	protectedUsers.Post("/logout", handlers.LogoutUser(authService))
}
