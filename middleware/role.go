package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
)

func HasRoles(roles ...entity.UserRoleName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("authenticated_user").(*entity.User)
		if !ok {
			return errors.UnauthorizedError("User not found", nil)
		}

		for _, role := range roles {
			if user.HasRole(role) {
				return c.Next()
			}
		}

		return errors.UnauthorizedError("User does not have the required role", nil)
	}
}
