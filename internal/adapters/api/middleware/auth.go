package middleware

import (
	"context"
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const AuthenticatedUserKey string = "auth_user"

var (
	errInvalidAuthFormat = "Invalid authorization header format"
	errTokenRequired     = "Access token is required"
)

func Protected(authService auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				response.Error(w, errors.UnauthorizedError(errTokenRequired, nil))
				return
			}

			if len(authorization) <= 7 || authorization[:7] != "Bearer " {
				response.Error(w, errors.UnauthorizedError(errInvalidAuthFormat, nil))
				return
			}

			accessToken := authorization[7:]
			user, err := authService.ValidateTokenAndGetUser(r.Context(), accessToken)
			if err != nil {
				response.Error(w, errors.UnauthorizedError(err.Error(), err))
				return
			}

			ctx := context.WithValue(r.Context(), AuthenticatedUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
