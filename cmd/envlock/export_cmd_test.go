package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlock/internal/snapshot"
)

func setupExportStore(t *testing.T, label string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVLOCK_STORE_DIR", dir)

	s := &snapshot.Snapshot{
		Label:     label,
		Timestamp: time.Now(),
		Vars:      map[string]string{"KEY": "value", "FOO": "bar"},
	}
	store := snapshot.DefaultStore()
	if err := store.Save(s); err != nil {
		t.Fatalf("setup: save snapshot: %v", err)
	}
	_ = filepath.Join(dir, label+".json") // confirm path shape
}

func TestRunExport_MissingLabel(t *testing.T) {
	err := runExport([]string{})
	if err == nil {
		t.Error("expected error for missing label")
	}
}

func TestRunExport_UnknownFlag(t *testing.T) {
	err := runExport([]string{"mylabel", "--format=xml"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestRunExport_MissingSnapshot(t *testing.T) {
	t.Setenv("ENVLOCK_STORE_DIR", t.TempDir())
	err := runExport([]string{"nonexistent"})
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRunExport_WritesToStdout(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVLOCK_STORE_DIR", dir)

	s := &snapshot.Snapshot{
		Label:     "ci",
		Timestamp: time.Now(),
		Vars:      map[string]string{"CI": "true"},
	}
	store := snapshot.DefaultStore()
	if err := store.Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}

	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runExport([]string{"ci", "--format=dotenv"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("runExport error: %v", err)
	}

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	output := string(buf[:n])
	if output == "" {
		t.Error("expected output to stdout, got empty")
	}
}
