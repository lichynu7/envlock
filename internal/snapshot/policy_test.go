package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

func makePolicy() snapshot.Policy {
	return snapshot.Policy{
		Name: "ci",
		Rules: []snapshot.PolicyRule{
			{Key: "CI", Required: true},
			{Key: "TOKEN", Required: true, MaxLen: 64},
			{Key: "OPTIONAL", Required: false, MaxLen: 10},
		},
	}
}

func makePolicySnapshot(vars map[string]string) snapshot.Snapshot {
	return snapshot.Snapshot{Vars: vars}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	p := makePolicy()
	s := makePolicySnapshot(map[string]string{
		"CI":    "true",
		"TOKEN": "abc123",
	})
	violations := snapshot.EnforcePolicy(s, p)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEnforcePolicy_MissingRequired(t *testing.T) {
	p := makePolicy()
	s := makePolicySnapshot(map[string]string{
		"CI": "true",
	})
	violations := snapshot.EnforcePolicy(s, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "TOKEN" {
		t.Errorf("expected violation for TOKEN, got %s", violations[0].Key)
	}
}

func TestEnforcePolicy_ValueTooLong(t *testing.T) {
	p := makePolicy()
	s := makePolicySnapshot(map[string]string{
		"CI":    "true",
		"TOKEN": "abc123",
		"OPTIONAL": "this-value-is-way-too-long",
	})
	violations := snapshot.EnforcePolicy(s, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "OPTIONAL" {
		t.Errorf("expected violation for OPTIONAL, got %s", violations[0].Key)
	}
}

func TestEnforcePolicy_OptionalAbsentIsOK(t *testing.T) {
	p := makePolicy()
	s := makePolicySnapshot(map[string]string{
		"CI":    "true",
		"TOKEN": "tok",
	})
	violations := snapshot.EnforcePolicy(s, p)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestSaveAndLoadPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	p := makePolicy()
	if err := snapshot.SavePolicy(path, p); err != nil {
		t.Fatalf("save policy: %v", err)
	}
	loaded, err := snapshot.LoadPolicy(path)
	if err != nil {
		t.Fatalf("load policy: %v", err)
	}
	if loaded.Name != p.Name {
		t.Errorf("name mismatch: got %s, want %s", loaded.Name, p.Name)
	}
	if len(loaded.Rules) != len(p.Rules) {
		t.Errorf("rule count mismatch: got %d, want %d", len(loaded.Rules), len(p.Rules))
	}
}

func TestLoadPolicy_Missing(t *testing.T) {
	_, err := snapshot.LoadPolicy(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
	_ = os.Stderr
}
