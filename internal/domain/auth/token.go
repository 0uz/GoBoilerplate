package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ouz/goboilerplate/pkg/auth"
	"github.com/ouz/goboilerplate/pkg/errors"
)

type Token struct {
	jwt.RegisteredClaims
	UserId     string          `json:"uid"`
	RawToken   string          `json:"token"`
	ClientType auth.ClientType `json:"clientType"`
	TokenType  auth.TokenType  `json:"tokenType"`
}

func validateTokenInput(userID string, tokenType auth.TokenType, jwtSecret string, clientType auth.ClientType, expiration time.Duration) error {
	if jwtSecret == "" {
		return errors.ValidationError("JWT secret cannot be empty", nil)
	}

	if clientType == "" {
		return errors.ValidationError("Client type cannot be empty", nil)
	}

	if userID == "" {
		return errors.ValidationError("User ID cannot be empty", nil)
	}

	if tokenType == "" {
		return errors.ValidationError("Token type cannot be empty", nil)
	}

	if expiration <= 0 {
		return errors.ValidationError("Expiration cannot be negative", nil)
	}

	return nil
}

func NewToken(jti, userID string, tokenType auth.TokenType, jwtSecret string, clientType auth.ClientType, expiration time.Duration) (Token, error) {
	if err := validateTokenInput(userID, tokenType, jwtSecret, clientType, expiration); err != nil {
		return Token{}, err
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := auth.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        jti,
		},
		UserId:     userID,
		TokenType:  tokenType,
		ClientType: clientType,
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return Token{}, errors.AuthError("Failed to generate token", err)
	}

	return Token{
		RegisteredClaims: claims.RegisteredClaims,
		UserId:           claims.UserId,
		RawToken:         tokenString,
		TokenType:        tokenType,
		ClientType:       clientType,
	}, nil
}

func ValidateToken(tokenString string, jwtSecret string) (*Token, error) {
	if jwtSecret == "" {
		return nil, errors.ValidationError("JWT secret cannot be empty", nil)
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.ExpiredTokenError("Token has expired", err)
		}
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := token.Claims.(*auth.TokenClaims)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	return &Token{
		RegisteredClaims: claims.RegisteredClaims,
		UserId:           claims.UserId,
		RawToken:         tokenString,
		ClientType:       claims.ClientType,
		TokenType:        claims.TokenType,
	}, nil
}

func GetTokenClaims(tokenString string, jwtSecret string) (*auth.TokenClaims, error) {
	if jwtSecret == "" {
		return nil, errors.ValidationError("JWT secret cannot be empty", nil)
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || token == nil {
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := token.Claims.(*auth.TokenClaims)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	return claims, nil
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Time)
}

func (t *Token) GetPrefix() string {
	if t.TokenType == auth.ACCESS_TOKEN {
		return fmt.Sprintf("uat:%s:%s", t.UserId, string(t.ClientType))
	}

	if t.TokenType == auth.REFRESH_TOKEN {
		return fmt.Sprintf("urt:%s:%s", t.UserId, string(t.ClientType))
	}

	return ""
}

func GeneratePrefix(tokenType auth.TokenType, userID string, clientType auth.ClientType) string {
	prefix := "uat"
	if tokenType == auth.REFRESH_TOKEN {
		prefix = "urt"
	}

	if clientType == "" {
		return fmt.Sprintf("%s:%s", prefix, userID)
	}

	return fmt.Sprintf("%s:%s:%s", prefix, userID, string(clientType))
}
