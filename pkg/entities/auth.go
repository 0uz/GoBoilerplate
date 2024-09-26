package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type CredentialType string
const (
	PASSWORD CredentialType = "PASSWORD"
)

type Credential struct {
	gorm.Model
	CredentialType CredentialType `gorm:"not null"`
	Hash           string         `gorm:"not null"`
	UserID         string         // Foreign key field
	User           User           `gorm:"foreignKey:UserID"`
}

type TokenType string
const (
	ACCESS_TOKEN  TokenType = "ACCESS_TOKEN"
	REFRESH_TOKEN TokenType = "REFRESH_TOKEN"
)

type Token struct {
	ID         uint      `gorm:"autoIncrement;primary_key"`
	Token      string    `gorm:"not null"`
	TokenType  TokenType `gorm:"not null"`
	Revoked    bool      `gorm:"default:false"`
	IpAddress  string    `gorm:"not null"`
	ClientType string    // Foreign key field
	Client     Client    `gorm:"foreignKey:ClientType;references:ClientType"`
	UserID     string    // Foreign key field
	User       User      `gorm:"foreignKey:UserID"`
	ExpiresAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ClientType string

const (
	IOS     ClientType = "IOS"
	MONITOR ClientType = "MONITOR"
)

type Client struct {
	ClientType   string     `gorm:"primary_key;not null"`
	ClientSecret string     `gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}
