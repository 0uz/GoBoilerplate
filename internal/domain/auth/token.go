package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type TokenType string

const (
	ACCESS_TOKEN  TokenType = "ACCESS_TOKEN"
	REFRESH_TOKEN TokenType = "REFRESH_TOKEN"
)

type Token struct {
	ID         string
	Token      string
	TokenType  TokenType
	Revoked    bool
	Client     Client
	ClientType string
	User       user.User
	UserID     string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"uid"`
}

// NewToken creates a new token instance with generated JWT
func NewToken(userID string, tokenType TokenType, clientType string, jwtSecret string, expiration time.Duration) (*Token, error) {
	jti := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        jti,
		},
		UserId: userID,
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.AuthError("Failed to generate token", err)
	}

	return &Token{
		ID:         jti,
		Token:      tokenString,
		TokenType:  tokenType,
		Revoked:    false,
		ClientType: clientType,
		UserID:     userID,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Validate validates the token and returns its claims
func (t *Token) Validate(jwtSecret string) (*TokenClaims, error) {
	if t.Revoked {
		return nil, errors.UnauthorizedError("Token is revoked", nil)
	}

	if t.IsExpired() {
		return nil, errors.UnauthorizedError("Token is expired", nil)
	}

	token, err := jwt.ParseWithClaims(t.Token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	return claims, nil
}

// Revoke marks the token as revoked
func (t *Token) Revoke() {
	t.Revoked = true
	t.UpdatedAt = time.Now()
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// GetCacheKey returns the cache key for this token based on its type and user
func (t *Token) GetCacheKey() string {
	prefix := "uat"
	if t.TokenType == REFRESH_TOKEN {
		prefix = "urt"
	}
	return fmt.Sprintf("%s:%s:%s", prefix, t.UserID, t.ClientType)
}
