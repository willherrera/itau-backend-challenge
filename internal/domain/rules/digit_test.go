package rules

import (
	"testing"
)

func TestDigitValidator(t *testing.T) {
	validator := NewDigitValidator()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with single digit",
			password: "a1",
			wantErr:  false,
		},
		{
			name:     "valid password with multiple digits",
			password: "abc123",
			wantErr:  false,
		},
		{
			name:     "valid password with only digits",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "invalid password without digits",
			password: "abcdef",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid password with special chars only",
			password: "!@#$%",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("DigitValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
