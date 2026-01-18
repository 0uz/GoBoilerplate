package auth

import (
	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

func NewTokenPair(userID string, clientType auth.ClientType, jwtConfig config.JWTConfig) (TokenPair, error) {
	jti := uuid.New().String()

	accessToken, err := NewToken(jti, userID, auth.ACCESS_TOKEN, jwtConfig.Secret, clientType, jwtConfig.AccessExpiration)
	if err != nil {
		return TokenPair{}, errors.AuthError("Failed to generate access token", err)
	}

	refreshToken, err := NewToken(jti, userID, auth.REFRESH_TOKEN, jwtConfig.Secret, clientType, jwtConfig.RefreshExpiration)
	if err != nil {
		return TokenPair{}, errors.AuthError("Failed to generate refresh token", err)
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
