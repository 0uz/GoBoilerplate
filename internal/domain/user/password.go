package user

import (
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Password is a value object that represents a user's password
type Password struct {
	hashed string
}

// NewPassword creates a new password instance with validation and hashing
func NewPassword(plaintext string) (*Password, error) {
	if err := validatePassword(plaintext); err != nil {
		return nil, err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.InternalError("Failed to hash password", err)
	}

	return &Password{hashed: string(hashedBytes)}, nil
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.ValidationError("Password must be at least 8 characters long", nil)
	}
	return nil
}

// Verify checks if the provided plaintext password matches the hashed password
func (p *Password) Verify(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashed), []byte(plaintext))
	return err == nil
}

// Hashed returns the hashed password string
func (p *Password) Hashed() string {
	return p.hashed
}
