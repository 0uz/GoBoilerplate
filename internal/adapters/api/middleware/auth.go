package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer"
)

func Protected(authService auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeader)
			if authHeader == "" {
				response.Error(w, errors.UnauthorizedError("missing authorization header", nil))
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 {
				response.Error(w, errors.UnauthorizedError("invalid authorization header format", nil))
				return
			}

			if !strings.EqualFold(parts[0], bearerPrefix) {
				response.Error(w, errors.UnauthorizedError("unsupported authorization type", nil))
				return
			}

			token := parts[1]
			user, err := authService.ValidateTokenAndGetUser(r.Context(), token)
			if err != nil {
				response.Error(w, errors.UnauthorizedError("invalid or expired token", err))
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), util.AuthenticatedUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
