package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlock/internal/snapshot"
)

func setupDiffStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVLOCK_DIR", dir)
	return dir
}

func saveDiffSnapshot(t *testing.T, dir, label string, vars map[string]string) {
	t.Helper()
	snap := &snapshot.Snapshot{
		Label: label,
		Vars:  vars,
	}
	store, err := snapshot.DefaultStore()
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("failed to save snapshot %q: %v", label, err)
	}
	_ = dir
}

func TestRunDiff_MissingLabels(t *testing.T) {
	setupDiffStore(t)
	err := runDiff([]string{})
	if err == nil {
		t.Fatal("expected error for missing labels")
	}
}

func TestRunDiff_OnlyOneLabel(t *testing.T) {
	setupDiffStore(t)
	err := runDiff([]string{"snap-a"})
	if err == nil {
		t.Fatal("expected error when only one label provided")
	}
}

func TestRunDiff_MissingSnapshot(t *testing.T) {
	setupDiffStore(t)
	err := runDiff([]string{"ghost-a", "ghost-b"})
	if err == nil {
		t.Fatal("expected error for non-existent snapshots")
	}
}

func TestRunDiff_UnknownFlag(t *testing.T) {
	setupDiffStore(t)
	err := runDiff([]string{"--notaflag"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRunDiff_Success(t *testing.T) {
	dir := setupDiffStore(t)

	saveDiffSnapshot(t, dir, "base", map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin",
	})
	saveDiffSnapshot(t, dir, "updated", map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/local/bin",
		"GOPATH": "/home/user/go",
	})

	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDiff([]string{"--color=false", "base", "updated"})

	w.Close()
	os.Stdout = old

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])
	_ = filepath.Join(dir, "base")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output == "" {
		t.Error("expected non-empty diff output")
	}
}
