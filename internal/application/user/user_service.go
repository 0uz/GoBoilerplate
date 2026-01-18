package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ouz/goauthboilerplate/pkg/cache"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	authDto "github.com/ouz/goauthboilerplate/internal/application/auth/dto"

	"github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/ouz/goauthboilerplate/pkg/log"
)

const (
	userCachePrefix = "user:%s"
	userCacheTTL    = 5 * time.Minute
)

type userService struct {
	userRepository user.UserRepository
	redisCache     cache.RedisCacheService
	tx             postgres.TransactionManager
	logger         *log.Logger
}

func NewUserService(logger *log.Logger, ur user.UserRepository, rc cache.RedisCacheService, tx postgres.TransactionManager) user.UserService {
	return &userService{
		userRepository: ur,
		redisCache:     rc,
		tx:             tx,
		logger:         logger,
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
		return errors.InternalError("Failed to check existing user", err)
	}

	if existingUser != nil {
		return errors.ConflictError("User already exists", nil)
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return errors.InternalError("Failed to create user", err)
	}

	s.logger.Info("User registered successfully, verification email will be sent", "userID", user.ID, "email", user.Email)
	// TODO: Send email to user
	return nil
}

func (s *userService) RegisterAnonymousUser(ctx context.Context) (*user.User, error) {
	user, err := user.NewAnonymousUser()
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, errors.InternalError("Failed to create anonymous user", err)
	}

	s.logger.Info("Anonymous user registered successfully", "user_id", user.ID)
	return user, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalError("Failed to find user by email", err)
	}
	return user, nil
}

func (s *userService) FindUserWithRoles(ctx context.Context, id string, fromCache bool) (*user.User, error) {
	if fromCache {
		var user = &user.User{}
		cacheKey := fmt.Sprintf(userCachePrefix, id)
		if found, err := s.redisCache.Get(ctx, cacheKey, "", user); err == nil && found {
			return user, nil
		}
	}

	user, err := s.userRepository.FindUserWithRoles(ctx, id)
	if err != nil {
		return nil, errors.InternalError("Failed to find user with roles", err)
	}

	if user != nil && fromCache {
		cacheKey := fmt.Sprintf(userCachePrefix, id)
		if err := s.redisCache.Set(ctx, cacheKey, "", userCacheTTL, user); err != nil {
			s.logger.Error("Failed to cache user", "error", err, "user_id", id)
		}
	}

	return user, nil
}

func (s *userService) ConfirmUser(ctx context.Context, confirmation string) error {
	userConfirmation, err := s.userRepository.FindConfirmationByID(ctx, confirmation)
	if err != nil {
		return errors.InternalError("Failed to find user confirmation", err)
	}

	if userConfirmation == nil {
		return errors.NotFoundError("User confirmation not found", nil)
	}

	if userConfirmation.User.ID == "" {
		return errors.InternalError("Invalid user confirmation data", nil)
	}

	userConfirmation.User.Confirm()

	err = s.tx.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		if err := s.userRepository.DeleteConfirmation(ctx, userConfirmation); err != nil {
			return errors.InternalError("Failed to delete user confirmation", err)
		}

		if err := s.userRepository.Update(ctx, &userConfirmation.User); err != nil {
			return errors.InternalError("Failed to confirm user", err)
		}

		cacheKey := fmt.Sprintf(userCachePrefix, userConfirmation.User.ID)
		if err := s.redisCache.Evict(ctx, cacheKey, ""); err != nil {
			s.logger.Error("Failed to invalidate user cache", "error", err, "userID", userConfirmation.User.ID)
		}
		return nil
	})

	if err != nil {
		return errors.InternalError("Failed to confirm user", err)
	}

	s.logger.Info("User confirmed successfully", "userID", userConfirmation.User.ID)
	return nil
}
