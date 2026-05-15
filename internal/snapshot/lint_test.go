package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeLintSnapshot(vars map[string]string) Snapshot {
	return Snapshot{
		Label:     "lint-test",
		Timestamp: time.Now(),
		Vars:      vars,
	}
}

func TestLint_NoWarnings(t *testing.T) {
	s := makeLintSnapshot(map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin:/bin",
	})
	result := Lint(s, DefaultLintOptions())
	if !result.IsClean() {
		t.Errorf("expected no warnings, got %d", len(result.Warnings))
	}
}

func TestLint_EmptyValue(t *testing.T) {
	s := makeLintSnapshot(map[string]string{
		"EMPTY_VAR": "",
	})
	opts := DefaultLintOptions()
	result := Lint(s, opts)

	found := false
	for _, w := range result.Warnings {
		if w.Name == "empty-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected empty-value warning")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	s := makeLintSnapshot(map[string]string{
		"myVar": "hello",
	})
	opts := DefaultLintOptions()
	result := Lint(s, opts)

	found := false
	for _, w := range result.Warnings {
		if w.Name == "lowercase-key" {
			found = true
		}
	}
	if !found {
		t.Error("expected lowercase-key warning")
	}
}

func TestLint_LongValue(t *testing.T) {
	s := makeLintSnapshot(map[string]string{
		"BIG": strings.Repeat("x", 600),
	})
	opts := DefaultLintOptions()
	result := Lint(s, opts)

	found := false
	for _, w := range result.Warnings {
		if w.Name == "long-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected long-value warning")
	}
}

func TestLint_LabelPreserved(t *testing.T) {
	s := makeLintSnapshot(map[string]string{"X": "1"})
	s.Label = "my-label"
	result := Lint(s, DefaultLintOptions())
	if result.Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", result.Label)
	}
}

func TestLint_DisabledRules(t *testing.T) {
	s := makeLintSnapshot(map[string]string{
		"lower": "",
	})
	opts := LintOptions{
		WarnOnEmpty:         false,
		WarnOnLowercaseKeys: false,
		WarnOnLongValues:    false,
	}
	result := Lint(s, opts)
	if !result.IsClean() {
		t.Errorf("expected clean result with all rules disabled, got %d warnings", len(result.Warnings))
	}
}
