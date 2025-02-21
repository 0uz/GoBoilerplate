package shared

import (
	"regexp"
	"strings"

	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type Email struct {
	Address string
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func NewEmail(address string) (Email, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return Email{}, errors.ValidationError("Email address cannot be empty", nil)
	}

	if !emailRegex.MatchString(address) {
		return Email{}, errors.ValidationError("Invalid email address format", nil)
	}

	// Split email into local and domain parts
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return Email{}, errors.ValidationError("Invalid email address format", nil)
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Additional validations
	if len(localPart) > 64 {
		return Email{}, errors.ValidationError("Local part of email cannot be longer than 64 characters", nil)
	}

	if len(domainPart) > 255 {
		return Email{}, errors.ValidationError("Domain part of email cannot be longer than 255 characters", nil)
	}

	// Check for consecutive dots
	if strings.Contains(address, "..") {
		return Email{}, errors.ValidationError("Email address cannot contain consecutive dots", nil)
	}

	return Email{Address: address}, nil
}
