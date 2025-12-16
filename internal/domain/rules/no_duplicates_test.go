package rules

import (
	"testing"
)

func TestNoDuplicatesValidator(t *testing.T) {
	validator := NewNoDuplicatesValidator()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with no duplicates",
			password: "AbTp9!fok",
			wantErr:  false,
		},
		{
			name:     "valid password with unique chars",
			password: "abc123!@#",
			wantErr:  false,
		},
		{
			name:     "invalid password with duplicate lowercase",
			password: "aa",
			wantErr:  true,
		},
		{
			name:     "invalid password with duplicate uppercase",
			password: "AA",
			wantErr:  true,
		},
		{
			name:     "invalid password with duplicate digit",
			password: "a1b1",
			wantErr:  true,
		},
		{
			name:     "invalid password with duplicate special char",
			password: "a!b!",
			wantErr:  true,
		},
		{
			name:     "invalid password from example - duplicate o",
			password: "AbTp9!foo",
			wantErr:  true,
		},
		{
			name:     "invalid password from example - duplicate A",
			password: "AbTp9!foA",
			wantErr:  true,
		},
		{
			name:     "invalid password with whitespace",
			password: "AbTp9 fok",
			wantErr:  true,
		},
		{
			name:     "invalid password with tab",
			password: "AbTp9\tfok",
			wantErr:  true,
		},
		{
			name:     "invalid password with newline",
			password: "AbTp9\nfok",
			wantErr:  true,
		},
		{
			name:     "invalid empty password",
			password: "",
			wantErr:  false, // Empty password has no duplicates
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoDuplicatesValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
