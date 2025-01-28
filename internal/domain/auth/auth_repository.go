package auth

import (
	"context"
)

type AuthRepository interface {
	FindClientBySecret(ctx context.Context, secret string) (*Client, error)
}
