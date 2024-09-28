package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	appError := errors.InternalError("An unexpected error occurred", err)

	if e, ok := err.(*errors.AppError); ok {
		appError = e
	}

	return c.Status(appError.Status).JSON(appError)
}
