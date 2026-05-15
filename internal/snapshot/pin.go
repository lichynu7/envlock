package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// PinnedVar represents a single pinned environment variable with metadata.
type PinnedVar struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	Comment   string    `json:"comment,omitempty"`
}

// PinStore manages pinned environment variables persisted to disk.
type PinStore struct {
	path string
}

// NewPinStore returns a PinStore backed by the given file path.
func NewPinStore(path string) *PinStore {
	return &PinStore{path: path}
}

func (s *PinStore) load() (map[string]PinnedVar, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]PinnedVar), nil
	}
	if err != nil {
		return nil, err
	}
	var pins map[string]PinnedVar
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, err
	}
	return pins, nil
}

func (s *PinStore) save(pins map[string]PinnedVar) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Pin adds or updates a pinned variable.
func (s *PinStore) Pin(key, value, comment string) error {
	pins, err := s.load()
	if err != nil {
		return err
	}
	pins[key] = PinnedVar{
		Key:      key,
		Value:    value,
		PinnedAt: time.Now().UTC(),
		Comment:  comment,
	}
	return s.save(pins)
}

// Unpin removes a pinned variable by key. Returns an error if the key is not pinned.
func (s *PinStore) Unpin(key string) error {
	pins, err := s.load()
	if err != nil {
		return err
	}
	if _, ok := pins[key]; !ok {
		return errors.New("pin not found: " + key)
	}
	delete(pins, key)
	return s.save(pins)
}

// List returns all currently pinned variables.
func (s *PinStore) List() ([]PinnedVar, error) {
	pins, err := s.load()
	if err != nil {
		return nil, err
	}
	result := make([]PinnedVar, 0, len(pins))
	for _, p := range pins {
		result = append(result, p)
	}
	return result, nil
}

// CheckDrift compares pinned variables against a snapshot's env map and returns
// keys whose current value differs from the pinned value.
func (s *PinStore) CheckDrift(env map[string]string) ([]string, error) {
	pins, err := s.load()
	if err != nil {
		return nil, err
	}
	var drifted []string
	for key, pin := range pins {
		if current, ok := env[key]; !ok || current != pin.Value {
			drifted = append(drifted, key)
		}
	}
	return drifted, nil
}
