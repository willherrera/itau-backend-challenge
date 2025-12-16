package rules

import (
	"testing"
)

func TestUppercaseValidator(t *testing.T) {
	validator := NewUppercaseValidator()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with single uppercase",
			password: "A1",
			wantErr:  false,
		},
		{
			name:     "valid password with multiple uppercase",
			password: "ABC123",
			wantErr:  false,
		},
		{
			name:     "valid password with only uppercase",
			password: "ABCDEF",
			wantErr:  false,
		},
		{
			name:     "invalid password without uppercase",
			password: "abc123",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid password with lowercase and digits only",
			password: "abc123!@#",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UppercaseValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
