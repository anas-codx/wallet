package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	Letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits       = "0123456789"
	Symbols      = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// PasswordConfig defines the configuration for the password generator (Builder Pattern).
type PasswordConfig struct {
	Length           int
	IncludeDigits    bool
	IncludeSymbols   bool
	IncludeUppercase bool
}

// GeneratePassword creates a secure random password based on the provided configuration.
func GeneratePassword(config PasswordConfig) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyz"
	if config.IncludeUppercase {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if config.IncludeDigits {
		charset += Digits
	}
	if config.IncludeSymbols {
		charset += Symbols
	}

	password := make([]byte, config.Length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}
