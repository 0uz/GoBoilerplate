package auth

import (
	"context"
)

type AuthRepository interface {
	SaveTokenPairs(ctx context.Context, tokens *[]Token) error
	RevokeAllOldTokens(ctx context.Context, userID string) error
	FindClientBySecret(ctx context.Context, secret string) (*Client, error)
	IsTokenRevoked(ctx context.Context, jti, userId string) (bool, error)
}
