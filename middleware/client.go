package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
)

func ClientSecret() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientSecret := entity.GetClientSecret(c)
		if clientSecret == "" {
			return errors.UnauthorizedError("Client secret is required", nil)
		}
		return c.Next()
	}
}