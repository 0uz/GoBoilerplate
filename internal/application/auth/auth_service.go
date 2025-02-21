package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	userAccessTokenPrefix  = "uat:%s:%s" // userID:clientType
	userRefreshTokenPrefix = "urt:%s:%s" // userID:clientType
	clientCachePrefix      = "client"
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

func (s *authService) GenerateToken(ctx context.Context, userId string) ([]auth.Token, error) {
	now := time.Now()

	accessTokenId := uuid.New().String()
	conf := config.Get().JWT
	accessToken, err := generateToken(conf.Secret, userId, accessTokenId, conf.AccessExpiration)
	if err != nil {
		return nil, errors.AuthError("Failed to generate access token", err)
	}
	refreshTokenId := uuid.New().String()
	refreshToken, err := generateToken(conf.Secret, userId, refreshTokenId, conf.RefreshExpiration)
	if err != nil {
		return nil, errors.InternalError("Failed to generate refresh token", err)
	}

	client, err := util.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.RevokeAllTokensByClient(ctx, userId, client.ClientType); err != nil {
		return nil, errors.InternalError("Failed to revoke old tokens", err)
	}

	clientType := string(client.ClientType)

	tokens := []auth.Token{
		createTokenEntity(accessTokenId, accessToken, auth.ACCESS_TOKEN, userId, clientType, now.Add(conf.AccessExpiration)),
		createTokenEntity(refreshTokenId, refreshToken, auth.REFRESH_TOKEN, userId, clientType, now.Add(conf.RefreshExpiration)),
	}

	if err := s.saveTokenPairs(ctx, userId, accessTokenId, refreshTokenId, clientType, conf); err != nil {
		return nil, errors.InternalError("Failed to save token pairs", err)
	}
	return tokens, nil
}

func generateToken(jwtSecret, userID string, jti string, expiration time.Duration) (string, error) {
	claims := auth.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			ID:        jti,
		},
		UserId: userID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
}

func (s *authService) saveTokenPairs(ctx context.Context, userId, accessTokenId, refreshTokenId, clientType string, conf config.JWTConfig) error {
	accessTokenKey := fmt.Sprintf(userAccessTokenPrefix, userId, clientType)
	if err := s.redisCache.Set(ctx, accessTokenKey, accessTokenId, conf.AccessExpiration, 0); err != nil {
		return errors.InternalError("Failed to save access token", err)
	}

	refreshTokenKey := fmt.Sprintf(userRefreshTokenPrefix, userId, clientType)
	if err := s.redisCache.Set(ctx, refreshTokenKey, refreshTokenId, conf.RefreshExpiration, 0); err != nil {
		return errors.InternalError("Failed to save refresh token", err)
	}

	return nil
}

func (s *authService) FindClientBySecretCached(ctx context.Context, clientSecret string) (auth.Client, error) {
	var cachedClient auth.Client
	if found, _ := s.redisCache.Get(ctx, clientCachePrefix, clientSecret, cachedClient); found {
		return cachedClient, nil
	}

	clientFromDB, err := s.authRepository.FindClientBySecret(ctx, clientSecret)
	if err != nil {
		return auth.Client{}, err
	}

	if err := s.redisCache.Set(ctx, clientCachePrefix, clientSecret, 1*time.Hour, clientFromDB); err != nil {
		s.logger.WithError(err).Error("Failed to cache client")
	}

	return *clientFromDB, nil
}

func createTokenEntity(id string, token string, tokenType auth.TokenType, userID, clientType string, expiresAt time.Time) auth.Token {
	return auth.Token{
		ID:         id,
		Token:      token,
		TokenType:  tokenType,
		Revoked:    false,
		ClientType: clientType,
		UserID:     userID,
		ExpiresAt:  expiresAt,
	}
}

func (s *authService) RevokeAllTokensByClient(ctx context.Context, userID string, clientType auth.ClientType) error {
	accessTokenKey := fmt.Sprintf(userAccessTokenPrefix, userID, string(clientType))
	if err := s.redisCache.EvictByPrefix(ctx, accessTokenKey); err != nil {
		return errors.InternalError("Failed to revoke access tokens", err)
	}

	refreshTokenKey := fmt.Sprintf(userRefreshTokenPrefix, userID, string(clientType))
	if err := s.redisCache.EvictByPrefix(ctx, refreshTokenKey); err != nil {
		return errors.InternalError("Failed to revoke refresh tokens", err)
	}

	return nil
}

func (s *authService) RevokeAllTokens(ctx context.Context, userID string) error {
	accessTokenKey := fmt.Sprintf(userAccessTokenPrefix, userID, "")
	if err := s.redisCache.EvictByPrefix(ctx, accessTokenKey); err != nil {
		return errors.InternalError("Failed to revoke access tokens", err)
	}

	refreshTokenKey := fmt.Sprintf(userRefreshTokenPrefix, userID, "")
	if err := s.redisCache.EvictByPrefix(ctx, refreshTokenKey); err != nil {
		return errors.InternalError("Failed to revoke refresh tokens", err)
	}

	return nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) ([]auth.Token, error) {
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	client, err := util.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	revoked, err := s.IsRefreshTokenRevoked(ctx, claims.ID, claims.UserId, string(client.ClientType))
	if err != nil {
		return nil, errors.InternalError("Failed to check if token is revoked", err)
	}

	if revoked {
		return nil, errors.UnauthorizedError("Token is revoked", nil)
	}

	user, err := s.userService.FindUserWithRoles(ctx, claims.UserId, true)
	if err != nil {
		return nil, errors.NotFoundError("User not found", err)
	}

	tokens, err := s.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &auth.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		config := config.Get().JWT
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := t.Claims.(*auth.TokenClaims)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	return claims, nil
}

func (s *authService) Login(ctx context.Context, email, password string) ([]auth.Token, error) {
	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalError("Failed to find user", err)
	}
	if user == nil {
		return nil, errors.NotFoundError("User not found", nil)
	}

	if !user.IsPasswordValid(password) {
		return nil, errors.UnauthorizedError("Invalid credentials", nil)
	}

	tokens, err := s.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *authService) LoginAnonymous(ctx context.Context, email string) ([]auth.Token, error) {
	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalError("Failed to find user", err)
	}
	if user == nil {
		return nil, errors.NotFoundError("User not found", nil)
	}

	tokens, err := s.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
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

	revoked, err := s.IsAccessTokenRevoked(ctx, claims.ID, claims.UserId, string(client.ClientType))
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

func (s *authService) IsAccessTokenRevoked(ctx context.Context, tokenID, userID, clientType string) (bool, error) {
	key := fmt.Sprintf(userAccessTokenPrefix, userID, clientType)
	exists, err := s.redisCache.Exists(ctx, key, tokenID)
	return !exists, err
}

func (s *authService) IsRefreshTokenRevoked(ctx context.Context, tokenID, userID, clientType string) (bool, error) {
	key := fmt.Sprintf(userRefreshTokenPrefix, userID, clientType)
	exists, err := s.redisCache.Exists(ctx, key, tokenID)
	return !exists, err
}
