package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CredentialType string

const (
	CredentialTypePassword CredentialType = "PASSWORD"
)

type Credential struct {
	ID             uint           `gorm:"primarykey"`
	CredentialType CredentialType `gorm:"not null"`
	Hash           string         `gorm:"not null"`
	User           User           `gorm:"foreignKey:UserID"`
	UserID         string         `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

func NewCredential(credentialType CredentialType, secret string) (*Credential, error) {
	password, err := NewPassword(secret)
	if err != nil {
		return nil, err
	}

	return &Credential{
		CredentialType: credentialType,
		Hash:           password.Hashed(),
	}, nil
}

func (c *Credential) IsPasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.Hash), []byte(password)) == nil
}
