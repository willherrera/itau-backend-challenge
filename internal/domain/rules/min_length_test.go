package rules

import (
	"testing"
)

func TestMinLengthValidator(t *testing.T) {
	validator := NewMinLengthValidator(9)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with exactly 9 characters",
			password: "123456789",
			wantErr:  false,
		},
		{
			name:     "valid password with more than 9 characters",
			password: "1234567890",
			wantErr:  false,
		},
		{
			name:     "invalid password with less than 9 characters",
			password: "12345678",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid password with 1 character",
			password: "a",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLengthValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
