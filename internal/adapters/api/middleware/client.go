package middleware

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

func HasClientSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientSecret := util.ExtractClientSecret(r)
		if clientSecret == "" {
			response.Error(w, errors.UnauthorizedError("Client secret is required", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}
