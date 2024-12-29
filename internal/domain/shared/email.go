package shared

import (
	"errors"
	"net/mail"
	"strings"
)

type Email struct {
	Address string
}

func NewEmail(address string) (Email, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return Email{}, errors.New("email address cannot be empty")
	}

	if _, err := mail.ParseAddress(address); err != nil {
		return Email{}, errors.New("invalid email address format")
	}

	return Email{Address: address}, nil
}
