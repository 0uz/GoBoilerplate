package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
	"github.com/ouz/gobackend/pkg/auth"
	entity "github.com/ouz/gobackend/pkg/entities"
)


func Protected(authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			return errors.UnauthorizedError("Invalid authorization header format", nil)
		}

		accessToken := strings.TrimPrefix(authorization, "Bearer ")

		if accessToken == "" {
			return errors.UnauthorizedError("Access token is required", nil)
		}

		user, err := authService.ValidateTokenAndGetUser(accessToken)
		if err != nil {
			return err
		}

		entity.SetAuthenticatedUser(c, user)

		return c.Next()
	}
}
