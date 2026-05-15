package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/envlock/internal/snapshot"
)

func tempAnnotationStore(t *testing.T) *snapshot.AnnotationStore {
	t.Helper()
	dir, err := os.MkdirTemp("", "envlock-annotate-*")
	if err != nil {
		t.Fatalf("tempAnnotationStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := snapshot.NewAnnotationStore(dir)
	if err != nil {
		t.Fatalf("NewAnnotationStore: %v", err)
	}
	return store
}

func TestAnnotationStore_SetAndGet(t *testing.T) {
	s := tempAnnotationStore(t)
	if err := s.Set("prod", "production baseline"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	a, err := s.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if a.Note != "production baseline" {
		t.Errorf("expected note %q, got %q", "production baseline", a.Note)
	}
	if a.Label != "prod" {
		t.Errorf("expected label %q, got %q", "prod", a.Label)
	}
}

func TestAnnotationStore_GetMissing(t *testing.T) {
	s := tempAnnotationStore(t)
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing annotation")
	}
}

func TestAnnotationStore_SetOverwritesNote(t *testing.T) {
	s := tempAnnotationStore(t)
	if err := s.Set("dev", "first note"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := s.Set("dev", "updated note"); err != nil {
		t.Fatalf("Set overwrite: %v", err)
	}
	a, err := s.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if a.Note != "updated note" {
		t.Errorf("expected updated note, got %q", a.Note)
	}
	if !a.UpdatedAt.After(a.CreatedAt) {
		t.Error("expected UpdatedAt to be after CreatedAt after overwrite")
	}
}

func TestAnnotationStore_Remove(t *testing.T) {
	s := tempAnnotationStore(t)
	if err := s.Set("staging", "staging note"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := s.Get("staging")
	if err == nil {
		t.Fatal("expected error after remove")
	}
}

func TestAnnotationStore_RemoveMissing(t *testing.T) {
	s := tempAnnotationStore(t)
	if err := s.Remove("ghost"); err == nil {
		t.Fatal("expected error removing nonexistent annotation")
	}
}

func TestAnnotationStore_List(t *testing.T) {
	s := tempAnnotationStore(t)
	for _, label := range []string{"a", "b", "c"} {
		if err := s.Set(label, label+"-note"); err != nil {
			t.Fatalf("Set %q: %v", label, err)
		}
	}
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 annotations, got %d", len(list))
	}
}

func TestAnnotationStore_ListEmpty(t *testing.T) {
	s := tempAnnotationStore(t)
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(list))
	}
}
