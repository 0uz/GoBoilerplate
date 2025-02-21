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
				response.Error(w, errors.UnauthorizedError("Client authentication required - missing client secret", nil))
				return
			}

			client, err := authService.FindClientBySecretCached(r.Context(), clientSecret)
			if err != nil {
				response.Error(w, errors.UnauthorizedError("Invalid client credentials", err))
				return
			}

			if client.DeletedAt != nil {
				response.Error(w, errors.ForbiddenError("Client is disabled", nil))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, util.ClientKey, client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
