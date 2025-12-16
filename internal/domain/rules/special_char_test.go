package rules

import (
	"testing"
)

func TestSpecialCharValidator(t *testing.T) {
	validator := NewSpecialCharValidator("!@#$%^&*()-+")

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with exclamation mark",
			password: "abc!",
			wantErr:  false,
		},
		{
			name:     "valid password with at symbol",
			password: "abc@123",
			wantErr:  false,
		},
		{
			name:     "valid password with multiple special chars",
			password: "abc!@#",
			wantErr:  false,
		},
		{
			name:     "valid password with all allowed special chars",
			password: "!@#$%^&*()-+",
			wantErr:  false,
		},
		{
			name:     "invalid password without special chars",
			password: "abc123ABC",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid password with non-allowed special char",
			password: "abc123_",
			wantErr:  true,
		},
		{
			name:     "invalid password with dot",
			password: "abc.123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecialCharValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
