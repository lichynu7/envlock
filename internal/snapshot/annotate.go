package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Annotation holds a user-defined note attached to a snapshot label.
type Annotation struct {
	Label     string    `json:"label"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnnotationStore manages annotations stored on disk.
type AnnotationStore struct {
	dir string
}

// NewAnnotationStore returns an AnnotationStore rooted at dir.
func NewAnnotationStore(dir string) (*AnnotationStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("annotate: mkdir: %w", err)
	}
	return &AnnotationStore{dir: dir}, nil
}

func (s *AnnotationStore) path(label string) string {
	return filepath.Join(s.dir, label+".annotation.json")
}

// Set creates or replaces the annotation for label.
func (s *AnnotationStore) Set(label, note string) error {
	now := time.Now().UTC()
	a := Annotation{Label: label, Note: note, CreatedAt: now, UpdatedAt: now}
	if existing, err := s.Get(label); err == nil {
		a.CreatedAt = existing.CreatedAt
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	return os.WriteFile(s.path(label), data, 0o600)
}

// Get retrieves the annotation for label.
func (s *AnnotationStore) Get(label string) (*Annotation, error) {
	data, err := os.ReadFile(s.path(label))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("annotate: no annotation for %q", label)
		}
		return nil, fmt.Errorf("annotate: read: %w", err)
	}
	var a Annotation
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("annotate: unmarshal: %w", err)
	}
	return &a, nil
}

// Remove deletes the annotation for label.
func (s *AnnotationStore) Remove(label string) error {
	err := os.Remove(s.path(label))
	if os.IsNotExist(err) {
		return fmt.Errorf("annotate: no annotation for %q", label)
	}
	return err
}

// List returns all stored annotations sorted by label.
func (s *AnnotationStore) List() ([]*Annotation, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.annotation.json"))
	if err != nil {
		return nil, fmt.Errorf("annotate: glob: %w", err)
	}
	var out []*Annotation
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		var a Annotation
		if err := json.Unmarshal(data, &a); err != nil {
			continue
		}
		out = append(out, &a)
	}
	return out, nil
}
