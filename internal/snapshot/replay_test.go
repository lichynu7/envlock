package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeReplaySnapshot(label string, ts time.Time, vars map[string]string) Snapshot {
	return Snapshot{Label: label, Timestamp: ts, Vars: vars}
}

func TestReplay_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)
	entries, err := Replay(store, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestReplay_OrderedByTimestamp(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)

	now := time.Now()
	s1 := makeReplaySnapshot("first", now.Add(-2*time.Minute), map[string]string{"A": "1"})
	s2 := makeReplaySnapshot("second", now.Add(-1*time.Minute), map[string]string{"A": "1", "B": "2"})
	s3 := makeReplaySnapshot("third", now, map[string]string{"A": "changed", "B": "2", "C": "3"})

	// save out of order intentionally
	_ = store.Save(s3)
	_ = store.Save(s1)
	_ = store.Save(s2)

	entries, err := Replay(store, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Label != "first" || entries[1].Label != "second" || entries[2].Label != "third" {
		t.Errorf("wrong order: %v %v %v", entries[0].Label, entries[1].Label, entries[2].Label)
	}
}

func TestReplay_DiffAttached(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)

	now := time.Now()
	s1 := makeReplaySnapshot("base", now.Add(-time.Minute), map[string]string{"A": "1"})
	s2 := makeReplaySnapshot("next", now, map[string]string{"A": "2", "B": "new"})
	_ = store.Save(s1)
	_ = store.Save(s2)

	entries, _ := Replay(store, DefaultReplayOptions())
	if entries[0].Diff != nil {
		t.Error("first entry should have no diff")
	}
	if entries[1].Diff == nil {
		t.Fatal("second entry should have a diff")
	}
	if len(entries[1].Diff.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(entries[1].Diff.Added))
	}
	if len(entries[1].Diff.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(entries[1].Diff.Changed))
	}
}

func TestReplay_NoDiff_WhenDisabled(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(dir)

	now := time.Now()
	s1 := makeReplaySnapshot("a", now.Add(-time.Minute), map[string]string{"X": "1"})
	s2 := makeReplaySnapshot("b", now, map[string]string{"X": "2"})
	_ = store.Save(s1)
	_ = store.Save(s2)

	opts := ReplayOptions{IncludeDiff: false}
	entries, _ := Replay(store, opts)
	for _, e := range entries {
		if e.Diff != nil {
			t.Errorf("step %d: expected no diff when disabled", e.Step)
		}
	}
}

func TestFormatReplay_Empty(t *testing.T) {
	out := FormatReplay(nil)
	if !strings.Contains(out, "no snapshots") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestFormatReplay_ContainsLabel(t *testing.T) {
	entry := ReplayEntry{Step: 1, Label: "mysnap", Timestamp: time.Now(), VarCount: 5}
	out := FormatReplay([]ReplayEntry{entry})
	if !strings.Contains(out, "mysnap") {
		t.Errorf("expected label in output, got: %q", out)
	}
	if !strings.Contains(out, "5 vars") {
		t.Errorf("expected var count in output, got: %q", out)
	}
}
