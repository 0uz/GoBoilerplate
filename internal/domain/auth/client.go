package auth

import (
	"time"

	"github.com/ouz/goboilerplate/pkg/auth"
)

type Client struct {
	ClientType   auth.ClientType `gorm:"primary_key;not null"`
	ClientSecret string          `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}
