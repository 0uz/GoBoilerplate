package user

import (
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashed string
}

func NewPassword(plaintext string) (*Password, error) {
	if len(plaintext) < 8 {
		return nil, errors.ValidationError("password must be at least 8 characters long", nil)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Password{hashed: string(hashedBytes)}, nil
}

func (p *Password) Verify(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashed), []byte(plaintext))
	return err == nil
}

func (p *Password) Hashed() string {
	return p.hashed
}
