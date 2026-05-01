package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileVault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	vaultPath := filepath.Join(tmpDir, "test.enc")
	password := "secret-pass"

	// 1. Test Creation
	store, err := NewStore(vaultPath, password)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 2. Test Set & Get
	err = store.Set("foo", "bar")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, ok := store.Get("foo")
	if !ok || val != "bar" {
		t.Errorf("Expected bar, got %s (ok: %v)", val, ok)
	}

	// 3. Test Persistence (Re-opening)
	store2, err := NewStore(vaultPath, password)
	if err != nil {
		t.Fatalf("Failed to re-open store: %v", err)
	}
	val2, ok := store2.Get("foo")
	if !ok || val2 != "bar" {
		t.Errorf("Persistence failed: expected bar, got %s", val2)
	}

	// 4. Test List
	store.Set("abc", "123")
	keys := store.List()
	if len(keys) != 2 || keys[0] != "abc" || keys[1] != "foo" {
		t.Errorf("List failed: got %v", keys)
	}

	// 5. Test Delete
	err = store.Delete("foo")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, ok = store.Get("foo")
	if ok {
		t.Error("Key still exists after delete")
	}

	// 6. Test Invalid Password
	_, err = NewStore(vaultPath, "wrong-pass")
	if err != ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}
