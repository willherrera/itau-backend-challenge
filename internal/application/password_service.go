package application

import (
	"github.com/willherrera/itau-backend-challenge/internal/domain"
)

type PasswordService struct {
	validators []domain.PasswordValidator
}

func NewPasswordService(validators []domain.PasswordValidator) *PasswordService {
	return &PasswordService{
		validators: validators,
	}
}

type ValidationResult struct {
	IsValid bool     `json:"isValid"`
	Errors  []string `json:"errors,omitempty"`
}

func (s *PasswordService) Validate(password string) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []string{},
	}

	for _, validator := range s.validators {
		if err := validator.Validate(password); err != nil {
			result.IsValid = false
			result.Errors = append(result.Errors, err.Error())
		}
	}

	return result
}
