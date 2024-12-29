package user

import (
	"context"
	"time"

	"github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	authDto "github.com/ouz/goauthboilerplate/internal/application/auth/dto"

	"github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type userService struct {
	userRepository user.UserRepository
	redisCache     redis.RedisCacheService
	tx             postgres.TransactionManager
}

func NewUserService(ur user.UserRepository, rc redis.RedisCacheService, tx postgres.TransactionManager) user.UserService {
	return &userService{
		userRepository: ur,
		redisCache:     rc,
		tx:             tx,
	}
}

func (s *userService) Register(ctx context.Context, request authDto.UserRegisterRequest) error {

	email, err := shared.NewEmail(request.Email)
	if err != nil {
		return err
	}

	user, err := user.NewUser(request.Username, request.Password, email)
	if err != nil {
		return err
	}

	existingUser, err := s.userRepository.FindNotVerifiedUser(ctx, user.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.ConflictError("user already exists", nil)
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	// TODO: Send email to user
	return nil
}

func (s *userService) RegisterAnonymousUser(ctx context.Context) (*user.User, error) {
	user, err := user.NewAnonymousUser()
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.userRepository.FindByEmail(ctx, email)
}

func (s *userService) FindUserWithRoles(ctx context.Context, id string, fromCache bool) (*user.User, error) {
	if fromCache {
		var user = &user.User{}
		if _, found := s.redisCache.Get(ctx, "user", id, user); found {
			return user, nil
		}
	}

	user, err := s.userRepository.FindUserWithRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	if fromCache {
		s.redisCache.Set(ctx, "user", id, 5*time.Minute, user)
	}

	return user, nil
}

func (s *userService) ConfirmUser(ctx context.Context, confirmation string) error {
	userConfirmation, err := s.userRepository.FindConfirmationByID(ctx, confirmation)
	if err != nil {
		return errors.InternalError("Failed to find user confirmation", err)
	}

	if userConfirmation == nil {
		return errors.NotFoundError("User confirmation not found", err)
	}

	userConfirmation.User.Confirm()

	err = s.tx.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		if err := s.userRepository.DeleteConfirmation(ctx, userConfirmation); err != nil {
			return errors.InternalError("Failed to delete user confirmation", err)
		}

		if err := s.userRepository.Update(ctx, &userConfirmation.User); err != nil {
			return errors.InternalError("Failed to confirm user", err)
		}
		return nil
	})

	if err != nil {
		return errors.InternalError("Failed to confirm user", err)
	}

	return nil
}
