package user

import (
	"time"

	"github.com/ouz/goauthboilerplate/pkg/errors"
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
	if err := validateCredentialType(credentialType); err != nil {
		return nil, err
	}

	password, err := NewPassword(secret)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Credential{
		CredentialType: credentialType,
		Hash:           password.Hashed(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func validateCredentialType(credType CredentialType) error {
	switch credType {
	case CredentialTypePassword:
		return nil
	default:
		return errors.ValidationError("Unsupported credential type", nil)
	}
}

func (c *Credential) IsPasswordValid(password string) bool {
	if c.CredentialType != CredentialTypePassword {
		return false
	}
	return (&Password{hashed: c.Hash}).Verify(password)
}
