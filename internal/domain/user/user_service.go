package user

import (
	"context"

	auth "github.com/ouz/goboilerplate/internal/application/auth/dto"
)

type UserService interface {
	Register(ctx context.Context, request auth.UserRegisterRequest) error
	RegisterAnonymousUser(ctx context.Context) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindUserWithRoles(ctx context.Context, id string, fromCache bool) (*User, error)
	ConfirmUser(ctx context.Context, confirmation string) error
}
