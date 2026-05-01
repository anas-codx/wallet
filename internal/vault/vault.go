package vault

import (
	"errors"
	"sort"

	"github.com/codingdestro/wallet-go/pkg/utils"
)

var (
	ErrInvalidPassword = errors.New("access denied: invalid password")
	ErrKeyNotFound      = errors.New("key not found in vault")
	ErrVaultNotFound    = errors.New("vault file not found")
)

type Vault struct {
	path     string
	password string
	data     map[string]string
}

// Open loads an existing vault or returns a new one if it doesn't exist.
func Open(path, password string) (*Vault, error) {
	v := &Vault{
		path:     path,
		password: password,
		data:     make(map[string]string),
	}

	if !utils.FileExists(path) {
		return v, nil
	}

	err := utils.LoadEncryptedJSON(path, password, &v.data)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return v, nil
}

// Save persists the vault data to disk.
func (v *Vault) Save() error {
	return utils.SaveEncryptedJSON(v.path, v.password, v.data)
}

// Set adds or updates a key-value pair.
func (v *Vault) Set(key, value string) error {
	v.data[key] = value
	return v.Save()
}

// Get retrieves a value by key.
func (v *Vault) Get(key string) (string, bool) {
	val, ok := v.data[key]
	return val, ok
}

// Delete removes a key from the vault.
func (v *Vault) Delete(key string) error {
	if _, ok := v.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(v.data, key)
	return v.Save()
}

// List returns a sorted list of all keys.
func (v *Vault) List() []string {
	keys := make([]string, 0, len(v.data))
	for k := range v.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Data returns a copy of the internal data map.
func (v *Vault) Data() map[string]string {
	copy := make(map[string]string)
	for k, v := range v.data {
		copy[k] = v
	}
	return copy
}
