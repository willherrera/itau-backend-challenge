package rules

import (
	"errors"
	"unicode"
)

type LowercaseValidator struct{}

func NewLowercaseValidator() *LowercaseValidator {
	return &LowercaseValidator{}
}

func (v *LowercaseValidator) Validate(password string) error {
	for _, char := range password {
		if unicode.IsLower(char) {
			return nil
		}
	}
	return ErrNoLowercase
}

var ErrNoLowercase = errors.New("password must contain at least one lowercase letter")
