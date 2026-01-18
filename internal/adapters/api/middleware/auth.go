package middleware

import (
	"context"
	"net/http"

	resp "github.com/ouz/goauthboilerplate/pkg/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	authorizationHeader = "Authorization"
)

func Protected(authService auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeader)
			if authHeader == "" {
				resp.Error(w, errors.UnauthorizedError("missing authorization header", nil))
				return
			}

			user, err := authService.ValidateTokenAndGetUser(r.Context(), authHeader)
			if err != nil {
				resp.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), util.AuthenticatedUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
