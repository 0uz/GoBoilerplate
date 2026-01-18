package middleware

import (
	"net/http"

	"slices"

	"github.com/ouz/goboilerplate/internal/adapters/api/util"
	"github.com/ouz/goboilerplate/internal/domain/user"
	"github.com/ouz/goboilerplate/pkg/errors"
	resp "github.com/ouz/goboilerplate/pkg/response"
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
