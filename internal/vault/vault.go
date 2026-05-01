package vault

import (
	"errors"
	"sort"

	"github.com/codingdestro/wallet-go/pkg/utils"
)

var (
	ErrInvalidPassword = errors.New("access denied: invalid password")
	ErrKeyNotFound      = errors.New("key not found in vault")
)

// SecretStore defines the contract for secret management (Dependency Inversion).
type SecretStore interface {
	Set(key, value string) error
	Get(key string) (string, bool)
	Delete(key string) error
	List() []string
}

// FileVault implements SecretStore using encrypted local files (Repository Pattern).
type FileVault struct {
	path     string
	password string
	data     map[string]string
}

// NewStore is a Factory Method to create a SecretStore.
func NewStore(path, password string) (SecretStore, error) {
	v := &FileVault{
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

func (v *FileVault) Save() error {
	return utils.SaveEncryptedJSON(v.path, v.password, v.data)
}

func (v *FileVault) Set(key, value string) error {
	v.data[key] = value
	return v.Save()
}

func (v *FileVault) Get(key string) (string, bool) {
	val, ok := v.data[key]
	return val, ok
}

func (v *FileVault) Delete(key string) error {
	if _, ok := v.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(v.data, key)
	return v.Save()
}

func (v *FileVault) List() []string {
	keys := make([]string, 0, len(v.data))
	for k := range v.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
