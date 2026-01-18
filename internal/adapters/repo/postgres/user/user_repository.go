package user

import (
	"context"

	"github.com/ouz/goboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goboilerplate/internal/domain/user"
	"github.com/ouz/goboilerplate/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	postgres.BaseRepository
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{BaseRepository: postgres.BaseRepository{
		DB: db,
	}}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := r.GetDB(ctx).Preload("Credentials").Where("email = ?", email).Where("enabled = ? AND verified = ?", true, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.InternalError("Failed to fetch user by email", err)
	}
	return &user, nil
}

func (r *userRepository) FindNotVerifiedUser(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := r.GetDB(ctx).Preload("Credentials").Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.InternalError("Failed to fetch not verified user by email", err)
	}
	return &user, nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	err := r.GetDB(ctx).Preload("Credentials").Where("enabled = ? AND verified = ?", true, true).Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("User not found", err)
		}
		return nil, errors.InternalError("Failed to fetch user by ID", err)
	}
	return &user, nil
}

func (r *userRepository) FindUserWithRoles(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	err := r.GetDB(ctx).WithContext(ctx).Preload("Roles").Where("id = ?", id).Where("enabled = ? AND verified = ?", true, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("User not found", err)
		}
		return nil, errors.InternalError("Failed to fetch user by ID", err)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {
	if err := r.GetDB(ctx).Create(user).Error; err != nil {
		return errors.InternalError("Failed to create user", err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	if err := r.GetDB(ctx).Model(user).Updates(user).Error; err != nil {
		return errors.InternalError("Failed to update user", err)
	}
	return nil
}

func (r *userRepository) FindConfirmationByID(ctx context.Context, id string) (*user.UserConfirmation, error) {
	var confirmation user.UserConfirmation
	err := r.GetDB(ctx).Where("id = ?", id).Preload("User").First(&confirmation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.InternalError("Failed to fetch user confirmation by ID", err)
	}
	return &confirmation, nil
}

func (r *userRepository) DeleteConfirmation(ctx context.Context, userConfirmation *user.UserConfirmation) error {
	if err := r.GetDB(ctx).Delete(userConfirmation).Error; err != nil {
		return errors.InternalError("Failed to delete user confirmation", err)
	}
	return nil
}

func (r *userRepository) CreateUserConfirmation(ctx context.Context, userConfirmation *user.UserConfirmation) error {
	if err := r.GetDB(ctx).Create(userConfirmation).Error; err != nil {
		return errors.InternalError("Failed to create user confirmation", err)
	}
	return nil
}
