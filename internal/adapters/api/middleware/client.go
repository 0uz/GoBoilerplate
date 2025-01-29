package middleware

import (
	"context"
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const ClientKey string = "client"

func HasClientSecret(authService auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientSecret := util.ExtractClientSecret(r)
			if clientSecret == "" {
				response.Error(w, errors.UnauthorizedError("Client secret is required", nil))
				return
			}

			client, err := authService.FindClientBySecretCached(r.Context(), clientSecret)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), ClientKey, client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
