package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/envlock/internal/snapshot"
)

func makeHistorySnapshot(label string, vars map[string]string, ts time.Time) snapshot.Snapshot {
	return snapshot.Snapshot{
		Label:     label,
		Timestamp: ts,
		Vars:      vars,
	}
}

func TestHistory_ReturnsEntriesSortedByTimestamp(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.DefaultStore(dir)

	now := time.Now()
	snaps := []snapshot.Snapshot{
		makeHistorySnapshot("oldest", map[string]string{"A": "1"}, now.Add(-2*time.Hour)),
		makeHistorySnapshot("newest", map[string]string{"A": "1", "B": "2"}, now),
		makeHistorySnapshot("middle", map[string]string{"A": "1", "B": "2", "C": "3"}, now.Add(-1*time.Hour)),
	}
	for _, s := range snaps {
		if err := store.Save(s); err != nil {
			t.Fatalf("Save(%q): %v", s.Label, err)
		}
	}

	entries, err := snapshot.History(store)
	if err != nil {
		t.Fatalf("History: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Label != "newest" {
		t.Errorf("expected newest first, got %q", entries[0].Label)
	}
	if entries[2].Label != "oldest" {
		t.Errorf("expected oldest last, got %q", entries[2].Label)
	}
}

func TestHistory_VarCountIsCorrect(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.DefaultStore(dir)

	snap := makeHistorySnapshot("test", map[string]string{"X": "1", "Y": "2", "Z": "3"}, time.Now())
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := snapshot.History(store)
	if err != nil {
		t.Fatalf("History: %v", err)
	}
	if entries[0].VarCount != 3 {
		t.Errorf("expected VarCount=3, got %d", entries[0].VarCount)
	}
}

func TestHistory_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.DefaultStore(dir)

	entries, err := snapshot.History(store)
	if err != nil {
		t.Fatalf("History: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestFormatHistory_Empty(t *testing.T) {
	out := snapshot.FormatHistory(nil)
	if !strings.Contains(out, "No snapshots") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestFormatHistory_ContainsLabels(t *testing.T) {
	entries := []snapshot.HistoryEntry{
		{Label: "prod-2024", Timestamp: time.Now(), VarCount: 12},
		{Label: "staging", Timestamp: time.Now().Add(-time.Hour), VarCount: 8},
	}
	out := snapshot.FormatHistory(entries)
	if !strings.Contains(out, "prod-2024") {
		t.Errorf("expected label prod-2024 in output")
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected label staging in output")
	}
	if !strings.Contains(out, "12") {
		t.Errorf("expected var count 12 in output")
	}
}
