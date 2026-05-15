package snapshot

import (
	"testing"
	"time"
)

func makePruneSnapshot(label string, ts time.Time) Snapshot {
	return Snapshot{
		Label:     label,
		Timestamp: ts,
		Vars:      map[string]string{"KEY": "value"},
	}
}

func TestPrune_KeepLast(t *testing.T) {
	store := tempStore(t)
	now := time.Now()

	for i, label := range []string{"snap-a", "snap-b", "snap-c"} {
		s := makePruneSnapshot(label, now.Add(time.Duration(i)*time.Minute))
		if err := store.Save(s); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	result, err := Prune(store, PruneOptions{KeepLast: 2})
	if err != nil {
		t.Fatalf("prune: %v", err)
	}

	if len(result.Removed) != 1 || result.Removed[0] != "snap-a" {
		t.Errorf("expected snap-a removed, got %v", result.Removed)
	}
	if len(result.Retained) != 2 {
		t.Errorf("expected 2 retained, got %d", len(result.Retained))
	}
}

func TestPrune_OlderThan(t *testing.T) {
	store := tempStore(t)
	now := time.Now()

	old := makePruneSnapshot("old-snap", now.Add(-48*time.Hour))
	recent := makePruneSnapshot("recent-snap", now.Add(-1*time.Hour))

	for _, s := range []Snapshot{old, recent} {
		if err := store.Save(s); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	cutoff := now.Add(-24 * time.Hour)
	result, err := Prune(store, PruneOptions{OlderThan: cutoff})
	if err != nil {
		t.Fatalf("prune: %v", err)
	}

	if len(result.Removed) != 1 || result.Removed[0] != "old-snap" {
		t.Errorf("expected old-snap removed, got %v", result.Removed)
	}
}

func TestPrune_DryRun(t *testing.T) {
	store := tempStore(t)
	now := time.Now()

	for i, label := range []string{"x", "y", "z"} {
		s := makePruneSnapshot(label, now.Add(time.Duration(i)*time.Minute))
		if err := store.Save(s); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	result, err := Prune(store, PruneOptions{KeepLast: 1, DryRun: true})
	if err != nil {
		t.Fatalf("prune dry run: %v", err)
	}

	if len(result.Removed) != 2 {
		t.Errorf("expected 2 reported removed, got %d", len(result.Removed))
	}

	// Snapshots should still exist.
	labels, _ := store.List()
	if len(labels) != 3 {
		t.Errorf("dry run should not delete; expected 3 snapshots, got %d", len(labels))
	}
}

func TestPrune_EmptyStore(t *testing.T) {
	store := tempStore(t)
	result, err := Prune(store, PruneOptions{KeepLast: 5})
	if err != nil {
		t.Fatalf("prune empty store: %v", err)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected nothing removed from empty store")
	}
}
