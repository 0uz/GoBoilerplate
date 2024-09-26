package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
	"github.com/ouz/gobackend/pkg/auth"
)

// Protected protect routes
func Protected(authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var access_token string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		}

		if access_token == "" {
			return errors.UnauthorizedError("Access token is required", nil)
		}

		user, err := authService.ValidateTokenAndGetUser(access_token)

		if err != nil {
			return errors.UnauthorizedError("Token is invalid", err)
		}

		c.Locals("authenticated_user", user)

		return c.Next()
	}
}
