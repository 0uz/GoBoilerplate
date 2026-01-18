package user

import (
	"time"

	"github.com/ouz/goboilerplate/pkg/errors"
	"gorm.io/gorm"
)

type UserRoleName string

const (
	UserRoleUser      UserRoleName = "USER"
	UserRoleAnonymous UserRoleName = "ANONYMOUS"
)

type UserRole struct {
	UserID    string       `gorm:"primaryKey;type:uuid"`
	Name      UserRoleName `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

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

func validateUserRole(name UserRoleName) error {
	switch name {
	case UserRoleUser, UserRoleAnonymous:
		return nil
	default:
		return errors.ValidationError("Unsupported role name", nil)
	}
}
