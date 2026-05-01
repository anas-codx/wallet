package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen    = 16
	keyLen     = 32 // AES-256
	iterations = 3
	memory     = 64 * 1024
	threads    = 4
)

// deriveKey derives a 32-byte key from a password and salt using Argon2id.
func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, iterations, uint32(memory), uint8(threads), keyLen)
}

// Encrypt encrypts data using AES-GCM with a password.
func Encrypt(password string, plaintext []byte) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Result format: [salt][nonce][ciphertext]
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	result := append(salt, append(nonce, ciphertext...)...)

	return result, nil
}

// Decrypt decrypts data using AES-GCM with a password.
func Decrypt(password string, data []byte) ([]byte, error) {
	if len(data) < saltLen {
		return nil, errors.New("data too short")
	}

	salt := data[:saltLen]
	remaining := data[saltLen:]

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(remaining) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := remaining[:nonceSize]
	ciphertext := remaining[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// SaveEncryptedJSON marshals data, encrypts it with a password, and saves it to a file.
func SaveEncryptedJSON(filename string, password string, data any) error {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(filename, ciphertext, 0o644)
}

// LoadEncryptedJSON reads an encrypted file, decrypts it with a password, and unmarshals it.
func LoadEncryptedJSON(filename string, password string, target any) error {
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	plaintext, err := Decrypt(password, ciphertext)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, target)
}
