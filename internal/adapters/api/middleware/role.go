package middleware

import (
	"net/http"

	"slices"

	resp "github.com/ouz/goauthboilerplate/pkg/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

func HasRoles(requiredRoles ...user.UserRoleName) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := util.GetAuthenticatedUser(r)
			if err != nil {
				resp.Error(w, err)
				return
			}

			hasRequiredRole := slices.ContainsFunc(requiredRoles, user.HasRole)

			if !hasRequiredRole {
				resp.Error(w, errors.ForbiddenError("Insufficient permissions", nil))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
