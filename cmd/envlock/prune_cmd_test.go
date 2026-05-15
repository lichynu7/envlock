package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envlock/internal/snapshot"
)

func setupPruneStore(t *testing.T) *snapshot.Store {
	t.Helper()
	store, err := snapshot.DefaultStore(t.TempDir())
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return store
}

func savePruneSnapshot(t *testing.T, store *snapshot.Store, label string, ts time.Time) {
	t.Helper()
	s := snapshot.Snapshot{
		Label:     label,
		Timestamp: ts,
		Vars:      map[string]string{"K": "v"},
	}
	if err := store.Save(s); err != nil {
		t.Fatalf("save snapshot %q: %v", label, err)
	}
}

func TestRunPrune_MissingFlags(t *testing.T) {
	store := setupPruneStore(t)
	var buf bytes.Buffer
	err := runPrune([]string{}, store, &buf)
	if err == nil || !strings.Contains(err.Error(), "at least one of") {
		t.Errorf("expected missing-flags error, got %v", err)
	}
}

func TestRunPrune_KeepLast(t *testing.T) {
	store := setupPruneStore(t)
	now := time.Now()
	savePruneSnapshot(t, store, "old", now.Add(-2*time.Minute))
	savePruneSnapshot(t, store, "mid", now.Add(-1*time.Minute))
	savePruneSnapshot(t, store, "new", now)

	var buf bytes.Buffer
	if err := runPrune([]string{"--keep-last=2"}, store, &buf); err != nil {
		t.Fatalf("runPrune: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "pruned: old") {
		t.Errorf("expected 'old' pruned in output, got: %s", out)
	}
	if !strings.Contains(out, "1 snapshot(s) pruned") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestRunPrune_DryRun(t *testing.T) {
	store := setupPruneStore(t)
	now := time.Now()
	savePruneSnapshot(t, store, "a", now.Add(-3*time.Minute))
	savePruneSnapshot(t, store, "b", now.Add(-2*time.Minute))
	savePruneSnapshot(t, store, "c", now)

	var buf bytes.Buffer
	if err := runPrune([]string{"--keep-last=1", "--dry-run"}, store, &buf); err != nil {
		t.Fatalf("runPrune dry-run: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[dry-run]") {
		t.Errorf("expected dry-run prefix, got: %s", out)
	}

	labels, _ := store.List()
	if len(labels) != 3 {
		t.Errorf("dry-run should not delete; got %d snapshots", len(labels))
	}
}

func TestRunPrune_NothingToPrune(t *testing.T) {
	store := setupPruneStore(t)
	savePruneSnapshot(t, store, "only", time.Now())

	var buf bytes.Buffer
	if err := runPrune([]string{"--keep-last=5"}, store, &buf); err != nil {
		t.Fatalf("runPrune: %v", err)
	}

	if !strings.Contains(buf.String(), "nothing to prune") {
		t.Errorf("expected 'nothing to prune', got: %s", buf.String())
	}
}
