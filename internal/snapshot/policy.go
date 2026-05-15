package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PolicyRule defines a single enforcement rule for environment variables.
type PolicyRule struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	MaxLen   int    `json:"max_len,omitempty"`
}

// Policy holds a named set of rules to validate snapshots against.
type Policy struct {
	Name  string       `json:"name"`
	Rules []PolicyRule `json:"rules"`
}

// PolicyViolation describes a single rule that was not satisfied.
type PolicyViolation struct {
	Key     string
	Message string
}

// EnforcePolicy checks a snapshot against a policy and returns any violations.
func EnforcePolicy(s Snapshot, p Policy) []PolicyViolation {
	var violations []PolicyViolation
	for _, rule := range p.Rules {
		val, exists := s.Vars[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, PolicyViolation{
				Key:     rule.Key,
				Message: "required variable is missing",
			})
			continue
		}
		if !exists {
			continue
		}
		if rule.MaxLen > 0 && len(val) > rule.MaxLen {
			violations = append(violations, PolicyViolation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value exceeds max length of %d", rule.MaxLen),
			})
		}
	}
	return violations
}

// LoadPolicy reads a JSON policy file from disk.
func LoadPolicy(path string) (Policy, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Policy{}, fmt.Errorf("read policy: %w", err)
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return Policy{}, fmt.Errorf("parse policy: %w", err)
	}
	return p, nil
}

// SavePolicy writes a policy to disk as JSON.
func SavePolicy(path string, p Policy) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal policy: %w", err)
	}
	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
