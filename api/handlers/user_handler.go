package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ouz/gobackend/api/presenter"
	"github.com/ouz/gobackend/errors"
	"github.com/ouz/gobackend/pkg/auth"
	entity "github.com/ouz/gobackend/pkg/entities"
	"github.com/ouz/gobackend/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AllUsers(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fetched, err := service.FindAll()
		if err != nil {
			return errors.InternalError("Failed to fetch users", err)
		}
		return c.JSON(presenter.UsersSuccessResponse(fetched))
	}
}

func RegisterUser(service user.Service, authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request *presenter.UserRegisterRequest

		if err := c.BodyParser(&request); err != nil {
			return errors.ValidationError("Invalid request body", err)
		}

		if err := presenter.ValidateStruct(request); err != nil {
			return err
		}

		user := &entity.User{
			Username: request.Email,
			Roles: []entity.UserRole{
				{
					Name: entity.UserRoleUser,
				},
			},
			Enabled:  false,
			Verified: false,
			Anonymous: false,
		}
		
		if err := service.RegisterUser(user); err != nil {
			return err
		}

		if err := authService.CreateCredentials(user, request.Password); err != nil {
			return errors.AuthError("Failed to create credentials", err)
		}

		return c.JSON(presenter.UserSuccessResponse(user))
	}
}

func RegisterAnonymousUser(service user.Service, authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &entity.User{
			Username: uuid.New().String(),
			Roles: []entity.UserRole{
				{
					Name: entity.UserRoleAnonymous,
				},
			},
			Enabled:  true,
			Verified: true,
			Anonymous: true,
		}

		if err := service.RegisterUser(user); err != nil {
			return err
		}

		password := uuid.New().String()

		if err := authService.CreateCredentials(user, password); err != nil {
			return errors.AuthError("Failed to create credentials", err)
		}

		return c.JSON(presenter.AnonymousUserSuccessResponse(user, password))
	}
}

func LoginUser(userService user.Service, authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.FormValue("username")
		password := c.FormValue("password")

		if email == "" || password == "" {
			return errors.ValidationError("Email and password are required", nil)
		}

		user, err := userService.FindByEmail(email)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NotFoundError("User not found", err)
			}
			return errors.InternalError("Failed to find user", err)
		}

		matched := false
		for _, credential := range user.Credentials {
			if credential.CredentialType == entity.PASSWORD {
				if err := bcrypt.CompareHashAndPassword([]byte(credential.Hash), []byte(password)); err == nil {
					matched = true
					break
				}
			}
		}

		if !matched {
			return errors.UnauthorizedError("Invalid credentials", nil)
		}

		tokens, err := authService.CreateToken(c, user)
		if err != nil {
			return errors.AuthError("Failed to create token", err)
		}

		return c.JSON(presenter.TokenResponse{
			AccessToken:  (*tokens)[0].Token,
			RefreshToken: (*tokens)[1].Token,
		})
	}
}

func RefreshAccessToken(service user.Service, authService auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &presenter.RefreshAccessTokenRequest{}

		if err := c.BodyParser(request); err != nil {
			return errors.ValidationError("Invalid request body", err)
		}

		if request.RefreshToken == "" {
			return errors.ValidationError("Refresh token is required", nil)
		}

		user, err := authService.ValidateTokenAndGetUser(request.RefreshToken)
		if err != nil {
			return errors.UnauthorizedError("Invalid refresh token", err)
		}

		tokens, err := authService.CreateToken(c, user)
		if err != nil {
			return errors.AuthError("Failed to create new tokens", err)
		}

		return c.JSON(presenter.TokenResponse{
			AccessToken:  (*tokens)[0].Token,
			RefreshToken: (*tokens)[1].Token,
		})
	}
}
