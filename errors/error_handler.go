package errors

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	appError := InternalError("An unexpected error occurred", err)

	if e, ok := err.(*AppError); ok {
		appError = e
	}

	return c.Status(appError.Status).JSON(appError)
}
