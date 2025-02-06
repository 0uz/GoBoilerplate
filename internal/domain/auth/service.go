package auth

import (
	"context"

	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

type AuthService interface {
	GenerateToken(ctx context.Context, userId string) ([]Token, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) ([]Token, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	Login(ctx context.Context, email, password string) ([]Token, error)
	LoginAnonymous(ctx context.Context, email string) ([]Token, error)
	Logout(ctx context.Context, userID string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateTokenAndGetUser(ctx context.Context, token string) (user.User, error)
	FindClientBySecretCached(ctx context.Context, clientSecret string) (Client, error)
}
