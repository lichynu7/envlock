package main

import (
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

func setupReplayStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func saveReplaySnapshot(t *testing.T, dir, label string, ts time.Time, vars map[string]string) {
	t.Helper()
	store, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	s := snapshot.Snapshot{Label: label, Timestamp: ts, Vars: vars}
	if err := store.Save(s); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestRunReplay_EmptyStore(t *testing.T) {
	dir := setupReplayStore(t)
	var buf strings.Builder
	code := runReplay([]string{"--store", dir}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "no snapshots") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestRunReplay_ListsSteps(t *testing.T) {
	dir := setupReplayStore(t)
	now := time.Now()
	saveReplaySnapshot(t, dir, "alpha", now.Add(-2*time.Minute), map[string]string{"X": "1"})
	saveReplaySnapshot(t, dir, "beta", now.Add(-time.Minute), map[string]string{"X": "1", "Y": "2"})

	var buf strings.Builder
	code := runReplay([]string{"--store", dir}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, buf.String())
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output: %q", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected 'beta' in output: %q", out)
	}
	if !strings.Contains(out, "step 1") || !strings.Contains(out, "step 2") {
		t.Errorf("expected step numbers in output: %q", out)
	}
}

func TestRunReplay_NoDiffFlag(t *testing.T) {
	dir := setupReplayStore(t)
	now := time.Now()
	saveReplaySnapshot(t, dir, "s1", now.Add(-time.Minute), map[string]string{"A": "1"})
	saveReplaySnapshot(t, dir, "s2", now, map[string]string{"A": "2"})

	var buf strings.Builder
	code := runReplay([]string{"--store", dir, "--no-diff"}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	// with no-diff the [+N -N ~N] summary should be absent
	if strings.Contains(buf.String(), "[+") {
		t.Errorf("did not expect diff summary with --no-diff: %q", buf.String())
	}
}

func TestRunReplay_InvalidFlag(t *testing.T) {
	var buf strings.Builder
	code := runReplay([]string{"--unknown-flag"}, &buf)
	if code != 1 {
		t.Errorf("expected exit 1 for unknown flag, got %d", code)
	}
}
