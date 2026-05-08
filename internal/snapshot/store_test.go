package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlock/internal/snapshot"
)

func tempStore(t *testing.T) *snapshot.Store {
	t.Helper()
	dir := t.TempDir()
	return &snapshot.Store{Dir: filepath.Join(dir, ".envlock")}
}

func TestStore_SaveAndLoad(t *testing.T) {
	st := tempStore(t)
	s := &snapshot.Snapshot{
		Label: "ci-baseline",
		Vars:  map[string]string{"CI": "true", "HOME": "/root"},
	}

	path, err := st.Save(s)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("saved file not found at %s: %v", path, err)
	}

	loaded, err := st.Load("ci-baseline")
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Vars["CI"] != "true" {
		t.Errorf("expected CI=true, got %q", loaded.Vars["CI"])
	}
}

func TestStore_LoadMissing(t *testing.T) {
	st := tempStore(t)
	_, err := st.Load("nonexistent")
	if err == nil {
		t.Error("expected error loading missing snapshot, got nil")
	}
}

func TestStore_List(t *testing.T) {
	st := tempStore(t)

	for _, label := range []string{"snap-a", "snap-b", "snap-c"} {
		_, err := st.Save(&snapshot.Snapshot{
			Label: label,
			Vars:  map[string]string{"X": "1"},
		})
		if err != nil {
			t.Fatalf("save %q error: %v", label, err)
		}
	}

	labels, err := st.List()
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(labels) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(labels))
	}
}

func TestStore_ListEmpty(t *testing.T) {
	st := tempStore(t)
	labels, err := st.List()
	if err != nil {
		t.Fatalf("unexpected error on empty store: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(labels))
	}
}

func TestStore_SaveUsesTimestampWhenNoLabel(t *testing.T) {
	st := tempStore(t)
	s := snapshot.Capture("")
	path, err := st.Save(s)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s: %v", path, err)
	}
}
