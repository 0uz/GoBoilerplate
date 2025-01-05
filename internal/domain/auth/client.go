package auth

import (
	"time"
)

type ClientType string

const (
	IOS     ClientType = "IOS"
	MONITOR ClientType = "MONITOR"
)

type Client struct {
	ClientType   string `gorm:"primary_key;not null"`
	ClientSecret string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}
