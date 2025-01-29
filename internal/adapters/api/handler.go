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

	clientSecretMiddleware := middleware.HasClientSecret(userAuthService)

	// Public routes with client secret
	userRouter.Handle("POST /login", clientSecretMiddleware(http.HandlerFunc(authHandler.LoginUser)))
	userRouter.Handle("POST /login/anonymous", clientSecretMiddleware(http.HandlerFunc(authHandler.LoginAnonymousUser)))
	userRouter.Handle("POST /register", clientSecretMiddleware(http.HandlerFunc(userHandler.RegisterUser)))
	userRouter.Handle("POST /register/anonymous", clientSecretMiddleware(http.HandlerFunc(userHandler.RegisterAnonymousUser)))

	userRouter.Handle("POST /token/refresh", clientSecretMiddleware(http.HandlerFunc(authHandler.RefreshAccessToken)))

	protectedUser := middleware.Chain(
		clientSecretMiddleware,
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("POST /logout", protectedUser(http.HandlerFunc(authHandler.LogoutUser)))
	userRouter.Handle("POST /logout/all", clientSecretMiddleware(http.HandlerFunc(authHandler.LogoutAll)))

	mainRouter.Handle("/auth/", http.StripPrefix("/auth", userRouter)) // Prefix all user routes with /user
}

func SetUpUserRoutes(mainRouter *http.ServeMux, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	userRouter := http.NewServeMux()

	protectedUser := middleware.Chain(
		middleware.HasClientSecret(userAuthService),
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("GET /me", protectedUser(http.HandlerFunc(userHandler.GetUser)))

	mainRouter.Handle("/users/", http.StripPrefix("/users", userRouter)) // Prefix all user routes with /user
}
