package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClientType string

const (
	IOS     ClientType = "IOS"
	ANDROID ClientType = "ANDROID"
	WEB     ClientType = "WEB"
)

type TokenType string

const (
	ACCESS_TOKEN  TokenType = "ACCESS_TOKEN"
	REFRESH_TOKEN TokenType = "REFRESH_TOKEN"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId     string     `json:"uid"`
	ClientType ClientType `json:"clientType"`
	TokenType  TokenType  `json:"tokenType"`
}
