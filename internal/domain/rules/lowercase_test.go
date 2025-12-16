package rules

import (
	"testing"
)

func TestLowercaseValidator(t *testing.T) {
	validator := NewLowercaseValidator()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with single lowercase",
			password: "a1",
			wantErr:  false,
		},
		{
			name:     "valid password with multiple lowercase",
			password: "abc123",
			wantErr:  false,
		},
		{
			name:     "valid password with only lowercase",
			password: "abcdef",
			wantErr:  false,
		},
		{
			name:     "invalid password without lowercase",
			password: "ABC123",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid password with uppercase and digits only",
			password: "ABC123!@#",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("LowercaseValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
