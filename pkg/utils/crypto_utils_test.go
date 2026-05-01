package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	password := "test-password"
	plaintext := []byte("hello world this is a secret message")

	// Test direct Encrypt/Decrypt
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(password, ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Expected %s, got %s", string(plaintext), string(decrypted))
	}

	// Test Decryption with wrong password
	_, err = Decrypt("wrong-password", ciphertext)
	if err == nil {
		t.Error("Decryption should have failed with wrong password")
	}
}

func TestEncryptedJSON(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	password := "secure-pass"
	data := TestData{Name: "Wallet", Value: 100}
	tmpFile := filepath.Join(os.TempDir(), "wallet_test.enc")
	
	// Clean up after test
	defer os.Remove(tmpFile)

	// Test Save
	err := SaveEncryptedJSON(tmpFile, password, data)
	if err != nil {
		t.Fatalf("SaveEncryptedJSON failed: %v", err)
	}

	// Test Load
	var loaded TestData
	err = LoadEncryptedJSON(tmpFile, password, &loaded)
	if err != nil {
		t.Fatalf("LoadEncryptedJSON failed: %v", err)
	}

	if loaded.Name != data.Name || loaded.Value != data.Value {
		t.Errorf("Loaded data mismatch. Got %+v, want %+v", loaded, data)
	}

	// Test Load with wrong password
	var wrongData TestData
	err = LoadEncryptedJSON(tmpFile, "wrong-pass", &wrongData)
	if err == nil {
		t.Error("LoadEncryptedJSON should have failed with wrong password")
	}
}
