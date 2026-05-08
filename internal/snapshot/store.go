package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultDir = ".envlock"

// Store manages persistence of snapshots on disk.
type Store struct {
	Dir string
}

// DefaultStore returns a Store using the default directory under the current
// working directory.
func DefaultStore() (*Store, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	return &Store{Dir: filepath.Join(cwd, defaultDir)}, nil
}

// Save writes a snapshot to disk as <label>.json (or timestamp if no label).
func (st *Store) Save(s *Snapshot) (string, error) {
	if err := os.MkdirAll(st.Dir, 0o755); err != nil {
		return "", fmt.Errorf("creating store dir: %w", err)
	}

	name := s.Label
	if name == "" {
		name = s.Timestamp.Format("20060102T150405Z")
	}
	path := filepath.Join(st.Dir, name+".json")

	data, err := s.Marshal()
	if err != nil {
		return "", fmt.Errorf("marshalling snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("writing snapshot file: %w", err)
	}
	return path, nil
}

// Load reads a snapshot from disk by label (filename without .json).
func (st *Store) Load(label string) (*Snapshot, error) {
	path := filepath.Join(st.Dir, label+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot %q: %w", label, err)
	}
	s, err := Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("parsing snapshot %q: %w", label, err)
	}
	return s, nil
}

// List returns the labels of all stored snapshots.
func (st *Store) List() ([]string, error) {
	entries, err := os.ReadDir(st.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading store dir: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			labels = append(labels, e.Name()[:len(e.Name())-5])
		}
	}
	return labels, nil
}
