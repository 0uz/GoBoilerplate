package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
)

func HasRoles(requiredRoles ...entity.UserRoleName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := entity.GetAuthenticatedUser(c)
		if user == nil {
			return errors.UnauthorizedError("Authentication required", nil)
		}

		for _, role := range requiredRoles {
			if user.HasRole(role) {
				return c.Next()
			}
		}

		return errors.UnauthorizedError("Insufficient permissions", nil)
	}
}
