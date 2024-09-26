package presenter

import (
	"github.com/gofiber/fiber/v2"
	v "github.com/ouz/gobackend/api/validator"
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
)

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AnonymousUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func UserSuccessResponse(user *entity.User) *fiber.Map {
	mappedUser := UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return &fiber.Map{
		"status": "success",
		"data":   mappedUser,
	}
}

func AnonymousUserSuccessResponse(user *entity.User, password string) AnonymousUserResponse {
	mappedUser := AnonymousUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Password: password,
	}
	return mappedUser
}

func UsersSuccessResponse(users []entity.User) *fiber.Map {
	mappedUsers := make([]UserResponse, len(users))
	for i, u := range users {
		mappedUsers[i] = UserResponse{
			ID:       u.ID,
			Username: u.Username,
		}
	}
	return &fiber.Map{
		"status": "success",
		"data":   mappedUsers,
	}
}

func ValidateStruct[T any](payload T) error {
	err := v.Validator.Struct(payload)
	if err != nil {
		return errors.ValidationError("Validation failed", err)
	}
	return nil
}
