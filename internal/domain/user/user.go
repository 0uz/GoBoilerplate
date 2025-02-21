package user

import (
	"time"

	"github.com/google/uuid"
	vo "github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/gorm"
)

// User represents the user aggregate root
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

// NewUser creates a new regular user with validation
func NewUser(username, password string, email vo.Email) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	now := time.Now()
	userID := uuid.New().String()

	role, err := NewUserRole(userID, UserRoleUser)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        userID,
		Username:  username,
		Email:     email.Address,
		Enabled:   false,
		Verified:  false,
		Anonymous: false,
		Roles:     []UserRole{*role},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := user.AddCredential(password); err != nil {
		return nil, err
	}

	if err := user.AddConfirmation(); err != nil {
		return nil, err
	}

	return user, nil
}

// validateUsername validates the username format
func validateUsername(username string) error {
	if username == "" {
		return errors.ValidationError("Username cannot be empty", nil)
	}
	if len(username) < 3 {
		return errors.ValidationError("Username must be at least 3 characters long", nil)
	}
	return nil
}

// NewAnonymousUser creates a new anonymous user
func NewAnonymousUser() (*User, error) {
	now := time.Now()
	userID := uuid.New().String()

	role, err := NewUserRole(userID, UserRoleAnonymous)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        userID,
		Username:  uuid.New().String(),
		Email:     uuid.New().String(),
		Roles:     []UserRole{*role},
		Enabled:   true,
		Verified:  true,
		Anonymous: true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddCredential adds a new credential to the user
func (u *User) AddCredential(password string) error {
	credential, err := NewCredential(CredentialTypePassword, password)
	if err != nil {
		return err
	}
	credential.UserID = u.ID
	u.Credentials = append(u.Credentials, *credential)
	u.UpdatedAt = time.Now()
	return nil
}

// AddConfirmation adds a new confirmation to the user
func (u *User) AddConfirmation() error {
	confirmation := UserConfirmation{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}
	u.Confirmations = append(u.Confirmations, confirmation)
	u.UpdatedAt = time.Now()
	return nil
}

// IsPasswordValid checks if the provided password is valid
func (u *User) IsPasswordValid(password string) bool {
	for _, credential := range u.Credentials {
		if credential.CredentialType == CredentialTypePassword && credential.IsPasswordValid(password) {
			return true
		}
	}
	return false
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(role UserRoleName) bool {
	for _, userRole := range u.Roles {
		if userRole.Name == role {
			return true
		}
	}
	return false
}

// Confirm confirms the user's account
func (u *User) Confirm() {
	u.Verified = true
	u.Enabled = true
	u.UpdatedAt = time.Now()
}
