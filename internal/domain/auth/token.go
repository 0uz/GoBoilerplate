package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type TokenType string

const (
	ACCESS_TOKEN  TokenType = "ACCESS_TOKEN"
	REFRESH_TOKEN TokenType = "REFRESH_TOKEN"
)

type Token struct {
	jwt.RegisteredClaims
	UserId     string     `json:"uid"`
	RawToken   string     `json:"token"`
	ClientType ClientType `json:"clientType"`
	TokenType  TokenType  `json:"tokenType"`
}

func validateTokenInput(userID string, tokenType TokenType, jwtSecret string, clientType ClientType, expiration time.Duration) error {
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

func NewToken(userID string, tokenType TokenType, jwtSecret string, clientType ClientType, expiration time.Duration) (Token, error) {
	if err := validateTokenInput(userID, tokenType, jwtSecret, clientType, expiration); err != nil {
		return Token{}, err
	}

	jti := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        jti,
		},
		UserId: userID,
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

// Validate validates the token and returns its claims
func ValidateToken(tokenString string, jwtSecret string) (*Token, error) {
	if jwtSecret == "" {
		return nil, errors.ValidationError("JWT secret cannot be empty", nil)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := token.Claims.(*Token)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	return claims, nil
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Time)
}

func (t *Token) GetPrefix() string {
	if t.TokenType == ACCESS_TOKEN {
		return fmt.Sprintf("uat:%s:%s", t.UserId, string(t.ClientType))
	}

	return fmt.Sprintf("urt:%s:%s", t.UserId, string(t.ClientType))
}

func GeneratePrefix(tokenType TokenType, userID string, clientType ClientType) string {
	prefix := "urt"
	if tokenType == ACCESS_TOKEN {
		prefix = "uat"
	}

	if clientType == "" {
		return fmt.Sprintf("%s:%s", prefix, userID)
	}

	return fmt.Sprintf("%s:%s:%s", prefix, userID, string(clientType))
}

func (t *Token) SetClient(client Client) {
	t.ClientType = client.ClientType
}
