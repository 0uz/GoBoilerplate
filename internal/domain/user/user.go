package user

import (
	"time"

	"github.com/google/uuid"
	vo "github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username      string
	Email         string
	Enabled       bool
	Verified      bool
	Anonymous     bool
	Roles         []UserRole         `gorm:"foreignKey:UserID"`
	Credentials   []Credential       `gorm:"foreignKey:UserID"`
	Confirmations []UserConfirmation `gorm:"foreignKey:UserID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func NewUser(username, password string, email vo.Email) (*User, error) {
	if username == "" {
		return nil, errors.ValidationError("Username is empty", nil)
	}

	credential, err := NewCredential(CredentialTypePassword, password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    email.Address,
		Credentials: []Credential{
			*credential,
		},
		Roles:         []UserRole{{Name: UserRoleUser}},
		Confirmations: []UserConfirmation{{}},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func NewAnonymousUser() (*User, error) {
	return &User{
		Username:  uuid.New().String(),
		Email:     uuid.New().String(),
		Roles:     []UserRole{{Name: UserRoleAnonymous}},
		Enabled:   true,
		Verified:  true,
		Anonymous: true,
	}, nil
}

func (u *User) IsPasswordValid(password string) bool {
	for _, credential := range u.Credentials {
		if credential.CredentialType == CredentialTypePassword && credential.IsPasswordValid(password) {
			return true
		}
	}
	return false
}

func (u *User) HasRole(role UserRoleName) bool {
	for _, userRole := range u.Roles {
		if userRole.Name == role {
			return true
		}
	}
	return false
}

func (u *User) Confirm() {
	u.Verified = true
	u.Enabled = true
}
