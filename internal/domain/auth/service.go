package auth

import (
	"context"

	"github.com/ouz/goboilerplate/internal/domain/user"
)

type AuthService interface {
	GenerateToken(ctx context.Context, userId string) (TokenPair, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (TokenPair, error)
	ValidateToken(ctx context.Context, token string) (*Token, error)
	Login(ctx context.Context, email, password string) (TokenPair, error)
	LoginAnonymous(ctx context.Context, email string) (TokenPair, error)
	Logout(ctx context.Context, userID string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateTokenAndGetUser(ctx context.Context, token string) (user.User, error)
	FindClientBySecretCached(ctx context.Context, clientSecret string) (Client, error)
}
