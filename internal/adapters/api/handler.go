package api

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

func SetUpUserRoutes(mainRouter *http.ServeMux, authHandler *AuthHandler, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET /confirm", userHandler.ConfirmUser)

	// Public routes with client secret
	userRouter.Handle("GET /login", middleware.HasClientSecret(http.HandlerFunc(authHandler.LoginUser)))
	userRouter.Handle("GET /anonymous/login", middleware.HasClientSecret(http.HandlerFunc(authHandler.LoginAnonymousUser)))
	userRouter.Handle("POST /register", middleware.HasClientSecret(http.HandlerFunc(userHandler.RegisterUser)))
	userRouter.Handle("POST /anonymous", middleware.HasClientSecret(http.HandlerFunc(userHandler.RegisterAnonymousUser)))
	userRouter.Handle("POST /refresh", middleware.HasClientSecret(http.HandlerFunc(authHandler.RefreshAccessToken)))

	protectedUser := middleware.Chain(
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("POST /logout", protectedUser(http.HandlerFunc(authHandler.LogoutUser)))

	mainRouter.Handle("/user/", http.StripPrefix("/user", userRouter)) // Prefix all user routes with /user
}
