package auth

import (
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

// NewTokenPair creates a new token pair for a user
func NewTokenPair(userID string, clientType string, jwtConfig config.JWTConfig) (*TokenPair, error) {
	accessToken, err := NewToken(userID, ACCESS_TOKEN, clientType, jwtConfig.Secret, jwtConfig.AccessExpiration)
	if err != nil {
		return nil, errors.AuthError("Failed to generate access token", err)
	}

	refreshToken, err := NewToken(userID, REFRESH_TOKEN, clientType, jwtConfig.Secret, jwtConfig.RefreshExpiration)
	if err != nil {
		return nil, errors.AuthError("Failed to generate refresh token", err)
	}

	return &TokenPair{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

// ToTokenSlice converts the token pair to a slice of tokens
func (tp *TokenPair) ToTokenSlice() []Token {
	return []Token{tp.AccessToken, tp.RefreshToken}
}
