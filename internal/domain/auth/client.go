package auth

import (
	"time"
)

type ClientType string

const (
	WEB     ClientType = "WEB"
	IOS     ClientType = "IOS"
	ANDROID ClientType = "ANDROID"
)

type Client struct {
	ClientType   ClientType `gorm:"primary_key;not null"`
	ClientSecret string     `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}
