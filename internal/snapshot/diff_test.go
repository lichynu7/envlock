package snapshot

import (
	"testing"
)

func makeSnapshot(vars map[string]string) *Snapshot {
	s := &Snapshot{}
	for k, v := range vars {
		s.Vars = append(s.Vars, EnvVar{Key: k, Value: v})
	}
	return s
}

func TestDiff_Added(t *testing.T) {
	base := makeSnapshot(map[string]string{"FOO": "bar"})
	target := makeSnapshot(map[string]string{"FOO": "bar", "NEW_VAR": "hello"})

	entries := Diff(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(entries))
	}
	if entries[0].Status != StatusAdded || entries[0].Key != "NEW_VAR" {
		t.Errorf("expected added NEW_VAR, got %+v", entries[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := makeSnapshot(map[string]string{"FOO": "bar", "OLD_VAR": "gone"})
	target := makeSnapshot(map[string]string{"FOO": "bar"})

	entries := Diff(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(entries))
	}
	if entries[0].Status != StatusRemoved || entries[0].Key != "OLD_VAR" {
		t.Errorf("expected removed OLD_VAR, got %+v", entries[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	base := makeSnapshot(map[string]string{"FOO": "bar"})
	target := makeSnapshot(map[string]string{"FOO": "baz"})

	entries := Diff(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Status != StatusChanged || e.Key != "FOO" || e.OldValue != "bar" || e.NewValue != "baz" {
		t.Errorf("unexpected diff entry: %+v", e)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	base := makeSnapshot(map[string]string{"FOO": "bar", "BAZ": "qux"})
	target := makeSnapshot(map[string]string{"FOO": "bar", "BAZ": "qux"})

	entries := Diff(base, target)
	if len(entries) != 0 {
		t.Errorf("expected no diff entries, got %d: %v", len(entries), entries)
	}
}

func TestDiffEntry_String(t *testing.T) {
	tests := []struct {
		entry    DiffEntry
		expected string
	}{
		{DiffEntry{Key: "X", NewValue: "1", Status: StatusAdded}, "+ X=1"},
		{DiffEntry{Key: "X", OldValue: "1", Status: StatusRemoved}, "- X=1"},
		{DiffEntry{Key: "X", OldValue: "a", NewValue: "b", Status: StatusChanged}, "~ X: \"a\" -> \"b\""},
	}
	for _, tt := range tests {
		got := tt.entry.String()
		if got != tt.expected {
			t.Errorf("String() = %q, want %q", got, tt.expected)
		}
	}
}
