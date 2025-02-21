package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	authDto "github.com/ouz/goauthboilerplate/internal/application/auth/dto"
	"github.com/sirupsen/logrus"

	"github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	userCachePrefix = "user:%s" // userID
	userCacheTTL    = 5 * time.Minute
)

type userService struct {
	userRepository user.UserRepository
	redisCache     redis.RedisCacheService
	tx             postgres.TransactionManager
	logger         *logrus.Logger
}

func NewUserService(logger *logrus.Logger, ur user.UserRepository, rc redis.RedisCacheService, tx postgres.TransactionManager) user.UserService {
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

	s.logger.WithFields(logrus.Fields{
		"userID": user.ID,
		"email":  user.Email,
	}).Info("User registered successfully, verification email will be sent")
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

	s.logger.WithField("user_id", user.ID).Info("Anonymous user registered successfully")
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
			s.logger.WithError(err).WithField("user_id", id).Error("Failed to cache user")
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

		// Invalidate user cache after confirmation
		cacheKey := fmt.Sprintf(userCachePrefix, userConfirmation.User.ID)
		if err := s.redisCache.Evict(ctx, cacheKey, ""); err != nil {
			s.logger.WithError(err).WithField("userID", userConfirmation.User.ID).Error("Failed to invalidate user cache")
		}

		return nil
	})

	if err != nil {
		return errors.InternalError("Failed to confirm user", err)
	}

	s.logger.WithField("userID", userConfirmation.User.ID).Info("User confirmed successfully")
	return nil
}
