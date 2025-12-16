package rules

import (
	"errors"
	"unicode"
)

type UppercaseValidator struct{}

func NewUppercaseValidator() *UppercaseValidator {
	return &UppercaseValidator{}
}

func (v *UppercaseValidator) Validate(password string) error {
	for _, char := range password {
		if unicode.IsUpper(char) {
			return nil
		}
	}
	return ErrNoUppercase
}

var ErrNoUppercase = errors.New("password must contain at least one uppercase letter")
