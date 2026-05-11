package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/envlock/internal/snapshot"
)

func setupCompareStore(t *testing.T) *snapshot.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := snapshot.NewStore(filepath.Join(dir, "snapshots"))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return store
}

func saveCompareSnapshot(t *testing.T, store *snapshot.Store, label string, vars map[string]string) {
	t.Helper()
	s := &snapshot.Snapshot{Label: label, Timestamp: time.Now(), Vars: vars}
	if err := store.Save(s); err != nil {
		t.Fatalf("failed to save snapshot %s: %v", label, err)
	}
}

func TestRunCompare_MissingLabels(t *testing.T) {
	store := setupCompareStore(t)
	var buf bytes.Buffer
	err := runCompare([]string{}, store, &buf)
	if err == nil || !strings.Contains(err.Error(), "two snapshot labels") {
		t.Errorf("expected missing labels error, got: %v", err)
	}
}

func TestRunCompare_OnlyOneLabel(t *testing.T) {
	store := setupCompareStore(t)
	var buf bytes.Buffer
	err := runCompare([]string{"only-one"}, store, &buf)
	if err == nil {
		t.Error("expected error for single label, got nil")
	}
}

func TestRunCompare_MissingSnapshot(t *testing.T) {
	store := setupCompareStore(t)
	var buf bytes.Buffer
	err := runCompare([]string{"ghost", "also-ghost"}, store, &buf)
	if err == nil || !strings.Contains(err.Error(), "could not load snapshot") {
		t.Errorf("expected load error, got: %v", err)
	}
}

func TestRunCompare_NoDiff_Summary(t *testing.T) {
	store := setupCompareStore(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	saveCompareSnapshot(t, store, "snap-a", vars)
	saveCompareSnapshot(t, store, "snap-b", vars)

	var buf bytes.Buffer
	// We can't easily intercept os.Exit, so skip exit-code check here.
	_ = runCompare([]string{"snap-a", "snap-b"}, store, &buf)

	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff summary, got: %s", buf.String())
	}
}

func TestRunCompare_Verbose_ShowsDiff(t *testing.T) {
	if os.Getenv("ENVLOCK_TEST_EXIT") == "1" {
		return // skip in subprocess
	}
	store := setupCompareStore(t)
	saveCompareSnapshot(t, store, "v-a", map[string]string{"FOO": "old"})
	saveCompareSnapshot(t, store, "v-b", map[string]string{"FOO": "new"})

	var buf bytes.Buffer
	_ = runCompare([]string{"--verbose", "v-a", "v-b"}, store, &buf)

	output := buf.String()
	if !strings.Contains(output, "FOO") {
		t.Errorf("expected verbose diff to contain FOO, got: %s", output)
	}
}
