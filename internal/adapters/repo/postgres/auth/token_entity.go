package postgres

import (
	"time"

	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

type TokenEntity struct {
	ID         string         `gorm:"type:uuid;primary_key"`
	TokenType  auth.TokenType `gorm:"not null"`
	Revoked    bool           `gorm:"default:false"`
	Client     auth.Client    `gorm:"foreignKey:ClientType;references:ClientType"`
	ClientType string         `gorm:"not null"`
	User       user.User      `gorm:"foreignKey:UserID"`
	UserID     string         `gorm:"not null"`
	ExpiresAt  time.Time      `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TokenEntity) TableName() string {
	return "tokens"
}

func FromDomain(token *auth.Token) TokenEntity {
	return TokenEntity{
		ID:         token.ID,
		TokenType:  token.TokenType,
		Revoked:    token.Revoked,
		Client:     token.Client,
		ClientType: token.ClientType,
		User:       token.User,
		UserID:     token.UserID,
		ExpiresAt:  token.ExpiresAt,
		CreatedAt:  token.CreatedAt,
		UpdatedAt:  token.UpdatedAt,
	}
}

func (t *TokenEntity) ToDomain() *auth.Token {
	return &auth.Token{
		ID:         t.ID,
		TokenType:  t.TokenType,
		Revoked:    t.Revoked,
		Client:     t.Client,
		ClientType: t.ClientType,
		User:       t.User,
		UserID:     t.UserID,
		ExpiresAt:  t.ExpiresAt,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
