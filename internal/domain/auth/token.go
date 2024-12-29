package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
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
