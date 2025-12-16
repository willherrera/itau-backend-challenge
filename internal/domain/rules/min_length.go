package rules

import (
	"fmt"
)

type MinLengthValidator struct {
	minLength int
}

func NewMinLengthValidator(minLength int) *MinLengthValidator {
	return &MinLengthValidator{
		minLength: minLength,
	}
}

func (v *MinLengthValidator) Validate(password string) error {
	if len(password) < v.minLength {
		return fmt.Errorf("password must have at least %d characters", v.minLength)
	}
	return nil
}
