package entities

import (
	"github.com/google/uuid"

	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"type:uuid;primary_key"`
	Username    string         `json:"username"`
	Enabled     bool           `json:"enabled"`
	Verified    bool           `json:"verified"`
	Anonymous   bool           `json:"anonymous"`
	Credentials []Credential   `gorm:"foreignKey:UserID"`
	Roles       []UserRole     `gorm:"foreignKey:UserID" json:"roles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

func (u *User) HasRole(role UserRoleName) bool {
	for _, userRole := range u.Roles {
		if userRole.Name == role {
			return true
		}
	}
	return false
}

type UserRoleName string
const (
	UserRoleUser      UserRoleName = "USER"
	UserRoleAnonymous UserRoleName = "ANONYMOUS"
)

type UserRole struct {
	UserID    string         `gorm:"primaryKey;type:uuid" json:"user_id"`
	Name      UserRoleName   `gorm:"primaryKey" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
