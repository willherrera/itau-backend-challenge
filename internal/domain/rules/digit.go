package rules

import (
	"errors"
	"unicode"
)

type DigitValidator struct{}

func NewDigitValidator() *DigitValidator {
	return &DigitValidator{}
}

func (v *DigitValidator) Validate(password string) error {
	for _, char := range password {
		if unicode.IsDigit(char) {
			return nil
		}
	}
	return ErrNoDigit
}

var ErrNoDigit = errors.New("password must contain at least one digit")
