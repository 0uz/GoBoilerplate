package auth

import (
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

func NewTokenPair(userID string, clientType ClientType, jwtConfig config.JWTConfig) (TokenPair, error) {
	accessToken, err := NewToken(userID, ACCESS_TOKEN, jwtConfig.Secret, clientType, jwtConfig.AccessExpiration)
	if err != nil {
		return TokenPair{}, errors.AuthError("Failed to generate access token", err)
	}

	refreshToken, err := NewToken(userID, REFRESH_TOKEN, jwtConfig.Secret, clientType, jwtConfig.RefreshExpiration)
	if err != nil {
		return TokenPair{}, errors.AuthError("Failed to generate refresh token", err)
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
