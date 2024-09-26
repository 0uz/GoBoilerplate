package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ouz/gobackend/errors"
)

func ErrorHandler(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			appErr := errors.GenericError("An unexpected error occurred", err)
			// logger.LogError(appErr)
			c.Status(appErr.Status).JSON(fiber.Map{
				"error": appErr,
			})
		}
	}()

	err := c.Next()

	if err != nil {
		var appErr *errors.AppError
		if e, ok := err.(*errors.AppError); ok {
			appErr = e
		} else {
			appErr = errors.GenericError(err.Error(), err)
		}
		// logger.LogError(appErr)
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error": appErr,
		})
	}

	return nil
}
