package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/cache/redis"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/sirupsen/logrus"
)

type authService struct {
	logger         *logrus.Logger
	authRepository auth.AuthRepository
	userService    user.UserService
	redisCache     redis.RedisCacheService
}

func NewAuthService(logger *logrus.Logger, ar auth.AuthRepository, us user.UserService, rc redis.RedisCacheService) auth.AuthService {
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

	client := util.GetClient(ctx)

	if client == nil {
		return nil, errors.InternalError("Failed to find client", err)
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
	err := s.redisCache.Set(ctx, "uat:"+userId+":"+clientType, accessTokenId, conf.AccessExpiration, 0)
	if err != nil {
		return err
	}
	return s.redisCache.Set(ctx, "urt:"+userId+":"+clientType, refreshTokenId, conf.RefreshExpiration, 0)
}

func (s *authService) FindClientBySecretCached(ctx context.Context, clientSecret string) (*auth.Client, error) {
	var cachedClient = &auth.Client{}
	if found, _ := s.redisCache.Get(ctx, "client", clientSecret, cachedClient); found {
		return cachedClient, nil
	}

	clientFromDB, err := s.authRepository.FindClientBySecret(ctx, clientSecret)
	if err != nil {
		return nil, err
	}

	s.redisCache.Set(ctx, "client", clientSecret, 1*time.Hour, clientFromDB)

	return clientFromDB, nil
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
	err := s.redisCache.EvictByPrefix(ctx, "uat:"+userID+":"+string(clientType))
	if err != nil {
		return err
	}

	err = s.redisCache.EvictByPrefix(ctx, "urt:"+userID+":"+string(clientType))
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) RevokeAllTokens(ctx context.Context, userID string) error {
	err := s.redisCache.EvictByPrefix(ctx, "uat:"+userID)
	if err != nil {
		return err
	}

	err = s.redisCache.EvictByPrefix(ctx, "urt:"+userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) ([]auth.Token, error) {
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	client := util.GetClient(ctx)
	if client == nil {
		return nil, errors.InternalError("Failed to get client", nil)
	}

	revoked, err := s.IsTokenRevoked(ctx, claims.ID, claims.UserId, string(client.ClientType))
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
	client := util.GetClient(ctx)
	if client == nil {
		return errors.InternalError("Failed to get client", nil)
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

func (s *authService) ValidateTokenAndGetUser(ctx context.Context, token string) (*user.User, error) {
	claims, err := s.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	client := util.GetClient(ctx)
	if client == nil {
		return nil, errors.InternalError("Failed to get client", nil)
	}

	revoked, err := s.IsTokenRevoked(ctx, claims.ID, claims.UserId, string(client.ClientType))
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

	return user, nil
}

func (s *authService) IsTokenRevoked(ctx context.Context, tokenID, userID, clientType string) (bool, error) {
	exists, err := s.redisCache.Exists(ctx, "uat:"+userID+":"+clientType, tokenID)
	return !exists, err
}
