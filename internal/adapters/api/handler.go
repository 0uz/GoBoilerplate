package api

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/middleware"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

func SetUpAuthRoutes(mainRouter *http.ServeMux, authHandler *AuthHandler, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	authRouter := http.NewServeMux()
	clientSecretMiddleware := middleware.HasClientSecret(userAuthService)

	// Public routes with client secret
	authRouter.Handle("POST /login", clientSecretMiddleware(http.HandlerFunc(authHandler.LoginUser)))
	authRouter.Handle("POST /login/anonymous", clientSecretMiddleware(http.HandlerFunc(authHandler.LoginAnonymousUser)))
	authRouter.Handle("POST /register", clientSecretMiddleware(http.HandlerFunc(userHandler.RegisterUser)))
	authRouter.Handle("POST /register/anonymous", clientSecretMiddleware(http.HandlerFunc(userHandler.RegisterAnonymousUser)))

	authRouter.Handle("POST /token/refresh", clientSecretMiddleware(http.HandlerFunc(authHandler.RefreshAccessToken)))

	protectedUser := middleware.Chain(
		clientSecretMiddleware,
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	authRouter.Handle("POST /logout", protectedUser(http.HandlerFunc(authHandler.LogoutUser)))
	authRouter.Handle("POST /logout/all", clientSecretMiddleware(http.HandlerFunc(authHandler.LogoutAll)))

	mainRouter.Handle("/auth/", http.StripPrefix("/auth", authRouter)) // Prefix all user routes with /user
}

func SetUpUserRoutes(mainRouter *http.ServeMux, userHandler *UserHandler, userAuthService auth.AuthService) {
	// Public routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET email/confirm", userHandler.ConfirmUser)
	userRouter.HandleFunc("GET email/confirm/resend", userHandler.ConfirmUser)

	protectedUser := middleware.Chain(
		middleware.HasClientSecret(userAuthService),
		middleware.Protected(userAuthService),
		middleware.HasRoles(user.UserRoleUser),
	)
	userRouter.Handle("GET /me", protectedUser(http.HandlerFunc(userHandler.GetUser)))

	mainRouter.Handle("/users/", http.StripPrefix("/users", userRouter)) // Prefix all user routes with /user
}
