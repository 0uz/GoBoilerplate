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

	user := &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email.Address,
		Roles:     []UserRole{{Name: UserRoleUser}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.AddCredential(password); err != nil {
		return nil, err
	}

	if err := user.AddConfirmation(); err != nil {
		return nil, err
	}

	return user, nil
}

func NewAnonymousUser() (*User, error) {
	return &User{
		ID:        uuid.New().String(),
		Username:  uuid.New().String(),
		Email:     uuid.New().String(),
		Roles:     []UserRole{{Name: UserRoleAnonymous}},
		Enabled:   true,
		Verified:  true,
		Anonymous: true,
		CreatedAt: time.Now(),
	}, nil
}

func (u *User) AddCredential(password string) error {
	credential, err := NewCredential(CredentialTypePassword, password)
	if err != nil {
		return err
	}
	u.Credentials = append(u.Credentials, *credential)
	return nil
}

func (u *User) AddConfirmation() error {
	confirmation := UserConfirmation{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}
	u.Confirmations = append(u.Confirmations, confirmation)
	return nil
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
	u.UpdatedAt = time.Now()
}
