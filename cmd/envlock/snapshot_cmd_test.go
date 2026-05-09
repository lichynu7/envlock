package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envlock/internal/snapshot"
)

// setupSnapshotStore creates a temporary directory and returns a store pointed at it,
// along with a cleanup function.
func setupSnapshotStore(t *testing.T) (*snapshot.Store, string) {
	t.Helper()
	dir := t.TempDir()
	store, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return store, dir
}

// TestRunSnapshot_CreatesSnapshot verifies that runSnapshot saves a snapshot
// with the given label to the store.
func TestRunSnapshot_CreatesSnapshot(t *testing.T) {
	store, _ := setupSnapshotStore(t)

	var out bytes.Buffer
	args := []string{"--label", "test-snap"}

	if err := runSnapshot(args, store, &out); err != nil {
		t.Fatalf("runSnapshot returned error: %v", err)
	}

	// Confirm the snapshot was persisted.
	snap, err := store.Load("test-snap")
	if err != nil {
		t.Fatalf("expected snapshot to be saved, got error: %v", err)
	}
	if snap.Label != "test-snap" {
		t.Errorf("expected label %q, got %q", "test-snap", snap.Label)
	}
}

// TestRunSnapshot_OutputConfirmation verifies that a confirmation message
// is written to stdout after a successful snapshot.
func TestRunSnapshot_OutputConfirmation(t *testing.T) {
	store, _ := setupSnapshotStore(t)

	var out bytes.Buffer
	args := []string{"--label", "ci-run-42"}

	if err := runSnapshot(args, store, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "ci-run-42") {
		t.Errorf("expected output to contain label %q, got: %q", "ci-run-42", output)
	}
}

// TestRunSnapshot_MissingLabel verifies that runSnapshot returns an error
// when no --label flag is provided.
func TestRunSnapshot_MissingLabel(t *testing.T) {
	store, _ := setupSnapshotStore(t)

	var out bytes.Buffer
	if err := runSnapshot([]string{}, store, &out); err == nil {
		t.Error("expected error when --label is missing, got nil")
	}
}

// TestRunSnapshot_CapturesCurrentEnv verifies that the captured snapshot
// contains at least one environment variable present in the current process.
func TestRunSnapshot_CapturesCurrentEnv(t *testing.T) {
	// Set a known env var so we can assert it was captured.
	const testKey = "ENVLOCK_TEST_VAR"
	const testVal = "hello_envlock"
	t.Setenv(testKey, testVal)

	store, _ := setupSnapshotStore(t)

	var out bytes.Buffer
	if err := runSnapshot([]string{"--label", "env-check"}, store, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := store.Load("env-check")
	if err != nil {
		t.Fatalf("failed to load snapshot: %v", err)
	}

	val, ok := snap.Vars[testKey]
	if !ok {
		t.Errorf("expected snapshot to contain %q", testKey)
	}
	if val != testVal {
		t.Errorf("expected %q=%q, got %q", testKey, testVal, val)
	}
}

// TestRunSnapshot_StoreDirectory verifies that the snapshot file is written
// to the expected store directory.
func TestRunSnapshot_StoreDirectory(t *testing.T) {
	store, dir := setupSnapshotStore(t)

	var out bytes.Buffer
	if err := runSnapshot([]string{"--label", "dir-check"}, store, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expect a file named after the label to exist in the store dir.
	expected := filepath.Join(dir, "dir-check.json")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("expected snapshot file at %q, but it does not exist", expected)
	}
}
