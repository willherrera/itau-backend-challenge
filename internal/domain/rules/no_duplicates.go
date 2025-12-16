package rules

import (
	"errors"
	"unicode"
)

type NoDuplicatesValidator struct{}

func NewNoDuplicatesValidator() *NoDuplicatesValidator {
	return &NoDuplicatesValidator{}
}

func (v *NoDuplicatesValidator) Validate(password string) error {
	seen := make(map[rune]bool)

	for _, char := range password {
		if unicode.IsSpace(char) {
			return ErrContainsWhitespace
		}

		if seen[char] {
			return ErrDuplicateChar
		}
		seen[char] = true
	}

	return nil
}

var ErrDuplicateChar = errors.New("password must not contain repeated characters")

var ErrContainsWhitespace = errors.New("password must not contain whitespace characters")
