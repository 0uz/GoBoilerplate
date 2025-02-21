package user

import (
	"time"

	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/gorm"
)

// UserRoleName represents the name of a user role
type UserRoleName string

const (
	UserRoleUser      UserRoleName = "USER"
	UserRoleAnonymous UserRoleName = "ANONYMOUS"
)

// UserRole represents a role assigned to a user
type UserRole struct {
	UserID    string       `gorm:"primaryKey;type:uuid"`
	Name      UserRoleName `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// NewUserRole creates a new user role with validation
func NewUserRole(userID string, name UserRoleName) (*UserRole, error) {
	if err := validateUserRole(name); err != nil {
		return nil, err
	}

	return &UserRole{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

// validateUserRole validates if the role name is supported
func validateUserRole(name UserRoleName) error {
	switch name {
	case UserRoleUser, UserRoleAnonymous:
		return nil
	default:
		return errors.ValidationError("Unsupported role name", nil)
	}
}
