package user

import (
	"github.com/ouz/goboilerplate/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashed string
}

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

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.ValidationError("Password must be at least 8 characters long", nil)
	}
	return nil
}

func (p *Password) Verify(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashed), []byte(plaintext))
	return err == nil
}

func (p *Password) Hashed() string {
	return p.hashed
}
