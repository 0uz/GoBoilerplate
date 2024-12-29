package user

import (
	"time"

	"gorm.io/gorm"
)

type UserConfirmation struct {
	ID        string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
