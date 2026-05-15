package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store persists snapshots on disk as JSON files.
type Store struct {
	dir string
}

// DefaultStore returns a Store rooted at dir, creating it if necessary.
func DefaultStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("store: mkdir %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(label string) string {
	return filepath.Join(s.dir, label+".json")
}

// Save writes snap to disk, overwriting any existing file with the same label.
func (s *Store) Save(snap Snapshot) error {
	data, err := Marshal(snap)
	if err != nil {
		return fmt.Errorf("store save: %w", err)
	}
	if err := os.WriteFile(s.path(snap.Label), data, 0o644); err != nil {
		return fmt.Errorf("store save: write %q: %w", snap.Label, err)
	}
	return nil
}

// Load reads and returns the snapshot with the given label.
func (s *Store) Load(label string) (Snapshot, error) {
	data, err := os.ReadFile(s.path(label))
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("snapshot %q not found", label)
		}
		return Snapshot{}, fmt.Errorf("store load: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("store load: unmarshal: %w", err)
	}
	return snap, nil
}

// List returns all snapshot labels stored in the directory.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("store list: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			labels = append(labels, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return labels, nil
}

// Delete removes the snapshot file for the given label.
func (s *Store) Delete(label string) error {
	if err := os.Remove(s.path(label)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot %q not found", label)
		}
		return fmt.Errorf("store delete %q: %w", label, err)
	}
	return nil
}
