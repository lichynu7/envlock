package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// BaselineEntry holds a pinned baseline snapshot label for a named environment.
type BaselineEntry struct {
	Environment string `json:"environment"`
	Label       string `json:"label"`
}

// BaselineStore manages baseline snapshot associations.
type BaselineStore struct {
	path string
}

// NewBaselineStore returns a BaselineStore backed by the given directory.
func NewBaselineStore(dir string) (*BaselineStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("baseline store: mkdir: %w", err)
	}
	return &BaselineStore{path: filepath.Join(dir, "baselines.json")}, nil
}

func (b *BaselineStore) load() (map[string]string, error) {
	data, err := os.ReadFile(b.path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline store: read: %w", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("baseline store: unmarshal: %w", err)
	}
	return m, nil
}

func (b *BaselineStore) save(m map[string]string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline store: marshal: %w", err)
	}
	return os.WriteFile(b.path, data, 0644)
}

// Set associates the given environment name with a snapshot label as its baseline.
func (b *BaselineStore) Set(environment, label string) error {
	m, err := b.load()
	if err != nil {
		return err
	}
	m[environment] = label
	return b.save(m)
}

// Get returns the baseline label for the given environment, or an error if not set.
func (b *BaselineStore) Get(environment string) (string, error) {
	m, err := b.load()
	if err != nil {
		return "", err
	}
	label, ok := m[environment]
	if !ok {
		return "", fmt.Errorf("baseline store: no baseline set for environment %q", environment)
	}
	return label, nil
}

// Remove deletes the baseline entry for the given environment.
func (b *BaselineStore) Remove(environment string) error {
	m, err := b.load()
	if err != nil {
		return err
	}
	if _, ok := m[environment]; !ok {
		return fmt.Errorf("baseline store: no baseline set for environment %q", environment)
	}
	delete(m, environment)
	return b.save(m)
}

// List returns all baseline entries.
func (b *BaselineStore) List() ([]BaselineEntry, error) {
	m, err := b.load()
	if err != nil {
		return nil, err
	}
	entries := make([]BaselineEntry, 0, len(m))
	for env, label := range m {
		entries = append(entries, BaselineEntry{Environment: env, Label: label})
	}
	return entries, nil
}
