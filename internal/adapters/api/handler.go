package api

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

func SetUpAuthRoutes(mainRouter *http.ServeMux, authHandler *AuthHandler, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET email/confirm", userHandler.ConfirmUser)
	userRouter.HandleFunc("GET email/confirm/resend", userHandler.ConfirmUser)

	// Public routes with client secret
	userRouter.Handle("POST /login", middleware.HasClientSecret(http.HandlerFunc(authHandler.LoginUser)))
	userRouter.Handle("POST /login/anonymous", middleware.HasClientSecret(http.HandlerFunc(authHandler.LoginAnonymousUser)))
	userRouter.Handle("POST /register", middleware.HasClientSecret(http.HandlerFunc(userHandler.RegisterUser)))
	userRouter.Handle("POST /register/anonymous", middleware.HasClientSecret(http.HandlerFunc(userHandler.RegisterAnonymousUser)))

	userRouter.Handle("POST /token/refresh", middleware.HasClientSecret(http.HandlerFunc(authHandler.RefreshAccessToken)))

	protectedUser := middleware.Chain(
		middleware.HasClientSecret,
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("POST auth/logout", protectedUser(http.HandlerFunc(authHandler.LogoutUser)))

	mainRouter.Handle("/auth/", http.StripPrefix("/auth", userRouter)) // Prefix all user routes with /user
}

func SetUpUserRoutes(mainRouter *http.ServeMux, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	userRouter := http.NewServeMux()

	protectedUser := middleware.Chain(
		middleware.HasClientSecret,
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("GET /me", protectedUser(http.HandlerFunc(userHandler.GetUser)))

	mainRouter.Handle("/user/", http.StripPrefix("/user", userRouter)) // Prefix all user routes with /user
}
