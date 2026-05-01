package ui

import (
	"testing"

	"github.com/codingdestro/wallet-go/internal/vault"
)

// MockStore implements vault.SecretStore for testing.
type MockStore struct {
	data map[string]string
}

func (m *MockStore) Set(k, v string) error { m.data[k] = v; return nil }
func (m *MockStore) Get(k string) (string, bool) { v, ok := m.data[k]; return v, ok }
func (m *MockStore) Delete(k string) error { delete(m.data, k); return nil }
func (m *MockStore) List() []string { return []string{"test-key"} }

// MockClipboard implements platform.Clipboard for testing.
type MockClipboard struct {
	content string
}

func (m *MockClipboard) Write(text string) error { m.content = text; return nil }

func TestUIModelTransitions(t *testing.T) {
	mockStore := &MockStore{data: make(map[string]string)}
	mockCB := &MockClipboard{}
	factory := func(pass string) (vault.SecretStore, error) {
		return mockStore, nil
	}

	model := NewModel(factory, mockCB)

	if model.State != StatePassword {
		t.Errorf("Expected StatePassword, got %v", model.State)
	}

	// Simulated events would go here (requires mocking tea.Msg types)
	// For now we test initialization and dependency setup.
	if model.StoreFactory == nil {
		t.Error("StoreFactory not initialized")
	}
}
