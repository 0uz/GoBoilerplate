package middleware

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

func HasRoles(requiredRoles ...user.UserRoleName) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := util.GetAuthenticatedUser(r)
			if err != nil {
				response.Error(w, err)
				return
			}

			for _, role := range requiredRoles {
				if user.HasRole(role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, errors.ForbiddenError("Insufficient permissions", nil))
		})
	}
}
