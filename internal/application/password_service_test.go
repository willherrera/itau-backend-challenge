package application

import (
	"testing"

	"github.com/willherrera/itau-backend-challenge/internal/domain"
	"github.com/willherrera/itau-backend-challenge/internal/domain/rules"
)

func TestPasswordService_Validate(t *testing.T) {
	service := NewPasswordService([]domain.PasswordValidator{
		rules.NewMinLengthValidator(9),
		rules.NewDigitValidator(),
		rules.NewLowercaseValidator(),
		rules.NewUppercaseValidator(),
		rules.NewSpecialCharValidator("!@#$%^&*()-+"),
		rules.NewNoDuplicatesValidator(),
	})

	tests := []struct {
		name      string
		password  string
		wantValid bool
		wantErrs  int // Expected number of errors
	}{
		{
			name:      "valid password from example",
			password:  "AbTp9!fok",
			wantValid: true,
			wantErrs:  0,
		},
		{
			name:      "invalid empty password",
			password:  "",
			wantValid: false,
			wantErrs:  5, // min_length, digit, lowercase, uppercase, special_char
		},
		{
			name:      "invalid password - too short",
			password:  "aa",
			wantValid: false,
			wantErrs:  5, // min_length, digit, uppercase, special_char, duplicates
		},
		{
			name:      "invalid password - no special chars",
			password:  "ab",
			wantValid: false,
			wantErrs:  4, // min_length, digit, uppercase, special_char
		},
		{
			name:      "invalid password - no digit",
			password:  "AAAbbbCc",
			wantValid: false,
			wantErrs:  4, // min_length, digit, special_char, duplicates
		},
		{
			name:      "invalid password - duplicate o",
			password:  "AbTp9!foo",
			wantValid: false,
			wantErrs:  1, // duplicates
		},
		{
			name:      "invalid password - duplicate A",
			password:  "AbTp9!foA",
			wantValid: false,
			wantErrs:  1, // duplicates
		},
		{
			name:      "invalid password - contains space",
			password:  "AbTp9 fok",
			wantValid: false,
			wantErrs:  2, // special_char, whitespace (from no_duplicates)
		},
		{
			name:      "invalid password - no uppercase",
			password:  "abtp9!fok",
			wantValid: false,
			wantErrs:  1, // uppercase
		},
		{
			name:      "invalid password - no lowercase",
			password:  "ABTP9!FOK",
			wantValid: false,
			wantErrs:  1, // lowercase
		},
		{
			name:      "valid password - all requirements met",
			password:  "Abc123!@#",
			wantValid: true,
			wantErrs:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.Validate(tt.password)

			if result.IsValid != tt.wantValid {
				t.Errorf("PasswordService.Validate() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrs {
				t.Errorf("PasswordService.Validate() got %d errors, want %d errors. Errors: %v",
					len(result.Errors), tt.wantErrs, result.Errors)
			}
		})
	}
}

func TestPasswordService_ValidateAllExamples(t *testing.T) {
	service := NewPasswordService([]domain.PasswordValidator{
		rules.NewMinLengthValidator(9),
		rules.NewDigitValidator(),
		rules.NewLowercaseValidator(),
		rules.NewUppercaseValidator(),
		rules.NewSpecialCharValidator("!@#$%^&*()-+"),
		rules.NewNoDuplicatesValidator(),
	})

	examples := []struct {
		password string
		expected bool
	}{
		{"", false},
		{"aa", false},
		{"ab", false},
		{"AAAbbbCc", false},
		{"AbTp9!foo", false},
		{"AbTp9!foA", false},
		{"AbTp9 fok", false},
		{"AbTp9!fok", true},
	}

	for _, ex := range examples {
		t.Run("example: "+ex.password, func(t *testing.T) {
			result := service.Validate(ex.password)
			if result.IsValid != ex.expected {
				t.Errorf("IsValid(%q) = %v, want %v. Errors: %v",
					ex.password, result.IsValid, ex.expected, result.Errors)
			}
		})
	}
}
