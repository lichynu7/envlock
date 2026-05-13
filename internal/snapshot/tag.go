package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// TagIndex maps tag names to snapshot labels.
type TagIndex map[string][]string

// TagStore manages a persistent tag index on disk.
type TagStore struct {
	path string
}

// NewTagStore returns a TagStore backed by the given directory.
func NewTagStore(dir string) *TagStore {
	return &TagStore{path: filepath.Join(dir, "tags.json")}
}

// Add associates a tag with a snapshot label.
func (ts *TagStore) Add(tag, label string) error {
	idx, err := ts.load()
	if err != nil {
		return err
	}
	for _, l := range idx[tag] {
		if l == label {
			return nil // already tagged
		}
	}
	idx[tag] = append(idx[tag], label)
	return ts.save(idx)
}

// Remove disassociates a tag from a snapshot label.
func (ts *TagStore) Remove(tag, label string) error {
	idx, err := ts.load()
	if err != nil {
		return err
	}
	updated := idx[tag][:0]
	for _, l := range idx[tag] {
		if l != label {
			updated = append(updated, l)
		}
	}
	if len(updated) == 0 {
		delete(idx, tag)
	} else {
		idx[tag] = updated
	}
	return ts.save(idx)
}

// Labels returns all snapshot labels associated with a tag, sorted.
func (ts *TagStore) Labels(tag string) ([]string, error) {
	idx, err := ts.load()
	if err != nil {
		return nil, err
	}
	labels := idx[tag]
	sort.Strings(labels)
	return labels, nil
}

// Tags returns all known tags, sorted.
func (ts *TagStore) Tags() ([]string, error) {
	idx, err := ts.load()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (ts *TagStore) load() (TagIndex, error) {
	idx := make(TagIndex)
	data, err := os.ReadFile(ts.path)
	if os.IsNotExist(err) {
		return idx, nil
	}
	if err != nil {
		return nil, fmt.Errorf("tag store read: %w", err)
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("tag store parse: %w", err)
	}
	return idx, nil
}

func (ts *TagStore) save(idx TagIndex) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("tag store marshal: %w", err)
	}
	return os.WriteFile(ts.path, data, 0644)
}
