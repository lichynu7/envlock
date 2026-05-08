package snapshot

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

// Snapshot represents a captured state of environment variables.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label,omitempty"`
	Vars      map[string]string `json:"vars"`
}

// Capture reads the current process environment and returns a Snapshot.
func Capture(label string) *Snapshot {
	envVars := os.Environ()
	vars := make(map[string]string, len(envVars))
	for _, entry := range envVars {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		} else if len(parts) == 1 {
			vars[parts[0]] = ""
		}
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Vars:      vars,
	}
}

// SortedKeys returns the environment variable keys in sorted order.
func (s *Snapshot) SortedKeys() []string {
	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Marshal serializes the snapshot to JSON bytes.
func (s *Snapshot) Marshal() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Unmarshal deserializes JSON bytes into a Snapshot.
func Unmarshal(data []byte) (*Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
