package rules

import (
	"errors"
	"strings"
)

type SpecialCharValidator struct {
	allowedChars string
}

func NewSpecialCharValidator(allowedChars string) *SpecialCharValidator {
	return &SpecialCharValidator{
		allowedChars: allowedChars,
	}
}

func (v *SpecialCharValidator) Validate(password string) error {
	for _, char := range password {
		if strings.ContainsRune(v.allowedChars, char) {
			return nil
		}
	}
	return ErrNoSpecialChar
}

var ErrNoSpecialChar = errors.New("password must contain at least one special character (!@#$%^&*()-+)")
