package user

import (
	"context"
)

type UserRepository interface {
	FindNotVerifiedUser(ctx context.Context, email string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	FindUserWithRoles(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	FindConfirmationByID(ctx context.Context, id string) (*UserConfirmation, error)
	DeleteConfirmation(ctx context.Context, userConfirmation *UserConfirmation) error
	CreateUserConfirmation(ctx context.Context, userConfirmation *UserConfirmation) error
}
