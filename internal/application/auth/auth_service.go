package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type authService struct {
	logger         *config.Logger
	authRepository auth.AuthRepository
	userService    user.UserService
	redisCache     redis.RedisCacheService
}

func NewAuthService(logger *config.Logger, ar auth.AuthRepository, us user.UserService, rc redis.RedisCacheService) auth.AuthService {
	return &authService{
		logger:         logger,
		authRepository: ar,
		userService:    us,
		redisCache:     rc,
	}
}

func (s *authService) GenerateToken(ctx context.Context, userId string) (auth.TokenPair, error) {
	client, err := util.GetClient(ctx)
	if err != nil {
		return auth.TokenPair{}, err
	}

	if err := s.RevokeAllTokensByClient(ctx, userId, client.ClientType); err != nil {
		return auth.TokenPair{}, errors.InternalError("Failed to revoke old tokens", err)
	}

	tokenPair, err := auth.NewTokenPair(userId, client.ClientType, config.Get().JWT)
	if err != nil {
		return auth.TokenPair{}, err
	}

	if err := s.saveTokenPair(ctx, tokenPair); err != nil {
		return auth.TokenPair{}, errors.InternalError("Failed to save token pair", err)
	}

	return tokenPair, nil
}

func (s *authService) saveTokenPair(ctx context.Context, tokenPair auth.TokenPair) error {
	if err := s.redisCache.Set(ctx, tokenPair.AccessToken.GetPrefix(), tokenPair.AccessToken.ID, config.Get().JWT.AccessExpiration, 0); err != nil {
		return errors.InternalError("Failed to save access token", err)
	}

	if err := s.redisCache.Set(ctx, tokenPair.RefreshToken.GetPrefix(), tokenPair.RefreshToken.ID, config.Get().JWT.RefreshExpiration, 0); err != nil {
		return errors.InternalError("Failed to save refresh token", err)
	}

	return nil
}

func (s *authService) FindClientBySecretCached(ctx context.Context, clientSecret string) (auth.Client, error) {
	var cachedClient auth.Client
	cacheKey := fmt.Sprintf("client:%s", clientSecret)

	if found, _ := s.redisCache.Get(ctx, cacheKey, clientSecret, &cachedClient); found {
		return cachedClient, nil
	}

	clientFromDB, err := s.authRepository.FindClientBySecret(ctx, clientSecret)
	if err != nil {
		return auth.Client{}, err
	}

	if err := s.redisCache.Set(ctx, cacheKey, clientSecret, 1*time.Hour, clientFromDB); err != nil {
		s.logger.Error("Failed to cache client", "error", err)
	}
	return *clientFromDB, nil
}

func (s *authService) RevokeAllTokensByClient(ctx context.Context, userID string, clientType auth.ClientType) error {
	accessTokenKey := auth.GeneratePrefix(auth.ACCESS_TOKEN, userID, clientType)
	if err := s.redisCache.EvictByPrefix(ctx, accessTokenKey); err != nil {
		return errors.InternalError("Failed to revoke access tokens", err)
	}

	refreshTokenKey := auth.GeneratePrefix(auth.REFRESH_TOKEN, userID, clientType)
	if err := s.redisCache.EvictByPrefix(ctx, refreshTokenKey); err != nil {
		return errors.InternalError("Failed to revoke refresh tokens", err)
	}

	return nil
}

func (s *authService) RevokeAllTokens(ctx context.Context, userID string) error {
	accessTokenKey := auth.GeneratePrefix(auth.ACCESS_TOKEN, userID, "")
	if err := s.redisCache.EvictByPrefix(ctx, accessTokenKey); err != nil {
		return errors.InternalError("Failed to revoke access tokens", err)
	}

	refreshTokenKey := auth.GeneratePrefix(auth.REFRESH_TOKEN, userID, "")
	if err := s.redisCache.EvictByPrefix(ctx, refreshTokenKey); err != nil {
		return errors.InternalError("Failed to revoke refresh tokens", err)
	}

	return nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (auth.TokenPair, error) {
	client, err := util.GetClient(ctx)
	if err != nil {
		return auth.TokenPair{}, err
	}

	claims, err := auth.ValidateToken(refreshToken, config.Get().JWT.Secret)
	if err != nil {
		return auth.TokenPair{}, err
	}

	claims.SetClient(client)
	claims.SetTokenType(auth.ACCESS_TOKEN)
	revoked, err := s.IsTokenRevoked(ctx, claims)
	if err != nil {
		return auth.TokenPair{}, errors.InternalError("Failed to check if token is revoked", err)
	}

	if revoked {
		return auth.TokenPair{}, errors.UnauthorizedError("Token is revoked", nil)
	}

	user, err := s.userService.FindUserWithRoles(ctx, claims.UserId, true)
	if err != nil {
		return auth.TokenPair{}, errors.NotFoundError("User not found", err)
	}

	return s.GenerateToken(ctx, user.ID)
}

func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (*auth.Token, error) {
	return auth.ValidateToken(tokenStr, config.Get().JWT.Secret)
}

func (s *authService) Login(ctx context.Context, email, password string) (auth.TokenPair, error) {
	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return auth.TokenPair{}, errors.InternalError("Failed to find user", err)
	}
	if user == nil {
		return auth.TokenPair{}, errors.NotFoundError("User not found", nil)
	}

	if !user.IsPasswordValid(password) {
		return auth.TokenPair{}, errors.UnauthorizedError("Invalid credentials", nil)
	}

	return s.GenerateToken(ctx, user.ID)
}

func (s *authService) LoginAnonymous(ctx context.Context, email string) (auth.TokenPair, error) {
	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return auth.TokenPair{}, errors.InternalError("Failed to find user", err)
	}
	if user == nil {
		return auth.TokenPair{}, errors.NotFoundError("User not found", nil)
	}

	return s.GenerateToken(ctx, user.ID)
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	client, err := util.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.RevokeAllTokensByClient(ctx, userID, client.ClientType); err != nil {
		return errors.InternalError("Failed to revoke old tokens", err)
	}
	return nil
}

func (s *authService) LogoutAll(ctx context.Context, userID string) error {
	if err := s.RevokeAllTokens(ctx, userID); err != nil {
		return errors.InternalError("Failed to revoke old tokens", err)
	}
	return nil
}

func (s *authService) ValidateTokenAndGetUser(ctx context.Context, token string) (user.User, error) {
	claims, err := s.ValidateToken(ctx, token)
	if err != nil {
		return user.User{}, err
	}

	client, err := util.GetClient(ctx)
	if err != nil {
		return user.User{}, err
	}

	claims.SetClient(client)
	claims.SetTokenType(auth.ACCESS_TOKEN)

	revoked, err := s.IsTokenRevoked(ctx, claims)
	if err != nil {
		return user.User{}, errors.InternalError("Failed to check if token is revoked", err)
	}

	if revoked {
		return user.User{}, errors.UnauthorizedError("Token is revoked", nil)
	}

	u, err := s.userService.FindUserWithRoles(ctx, claims.UserId, true)
	if err != nil {
		return user.User{}, errors.UnauthorizedError("Invalid user", err)
	}

	return *u, nil
}

func (s *authService) IsTokenRevoked(ctx context.Context, token *auth.Token) (bool, error) {
	key := token.GetPrefix()
	found, err := s.redisCache.Exists(ctx, key, token.ID)
	return !found, err
}
