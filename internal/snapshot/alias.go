package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// AliasStore maps human-friendly alias names to snapshot labels.
type AliasStore struct {
	mu   sync.RWMutex
	path string
	data map[string]string // alias -> label
}

// NewAliasStore creates or loads an alias store backed by the given file path.
func NewAliasStore(path string) (*AliasStore, error) {
	as := &AliasStore{
		path: path,
		data: make(map[string]string),
	}
	if err := as.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return as, nil
}

// Set associates an alias with a snapshot label, overwriting any existing mapping.
func (as *AliasStore) Set(alias, label string) error {
	if alias == "" || label == "" {
		return errors.New("alias and label must not be empty")
	}
	as.mu.Lock()
	defer as.mu.Unlock()
	as.data[alias] = label
	return as.save()
}

// Resolve returns the snapshot label for the given alias, or an error if not found.
func (as *AliasStore) Resolve(alias string) (string, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	label, ok := as.data[alias]
	if !ok {
		return "", errors.New("alias not found: " + alias)
	}
	return label, nil
}

// Remove deletes an alias mapping.
func (as *AliasStore) Remove(alias string) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	if _, ok := as.data[alias]; !ok {
		return errors.New("alias not found: " + alias)
	}
	delete(as.data, alias)
	return as.save()
}

// List returns all alias->label mappings.
func (as *AliasStore) List() map[string]string {
	as.mu.RLock()
	defer as.mu.RUnlock()
	copy := make(map[string]string, len(as.data))
	for k, v := range as.data {
		copy[k] = v
	}
	return copy
}

func (as *AliasStore) load() error {
	b, err := os.ReadFile(as.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &as.data)
}

func (as *AliasStore) save() error {
	if err := os.MkdirAll(filepath.Dir(as.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(as.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(as.path, b, 0o644)
}
