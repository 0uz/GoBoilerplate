package middleware

import (
	"context"
	"net/http"

	"github.com/ouz/goboilerplate/internal/adapters/api/util"
	"github.com/ouz/goboilerplate/internal/domain/auth"
	"github.com/ouz/goboilerplate/pkg/errors"
	resp "github.com/ouz/goboilerplate/pkg/response"
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
