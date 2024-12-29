package user

import (
	"time"

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
