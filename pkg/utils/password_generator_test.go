package utils

import "testing"

func TestGeneratePassword(t *testing.T) {
	configs := []PasswordConfig{
		{Length: 16, IncludeDigits: true, IncludeSymbols: true, IncludeUppercase: true},
		{Length: 8, IncludeDigits: false, IncludeSymbols: false, IncludeUppercase: false},
	}

	for _, config := range configs {
		pwd, err := GeneratePassword(config)
		if err != nil {
			t.Errorf("Failed to generate password: %v", err)
		}
		if len(pwd) != config.Length {
			t.Errorf("Expected length %d, got %d", config.Length, len(pwd))
		}
	}
}
