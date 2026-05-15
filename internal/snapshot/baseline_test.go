package snapshot

import (
	"os"
	"testing"
)

func tempBaselineStore(t *testing.T) *BaselineStore {
	t.Helper()
	dir, err := os.MkdirTemp("", "baseline-store-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := NewBaselineStore(dir)
	if err != nil {
		t.Fatalf("NewBaselineStore: %v", err)
	}
	return store
}

func TestBaselineStore_SetAndGet(t *testing.T) {
	store := tempBaselineStore(t)
	if err := store.Set("production", "snap-001"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	label, err := store.Get("production")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if label != "snap-001" {
		t.Errorf("expected snap-001, got %q", label)
	}
}

func TestBaselineStore_GetMissing(t *testing.T) {
	store := tempBaselineStore(t)
	_, err := store.Get("staging")
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}

func TestBaselineStore_SetOverwrites(t *testing.T) {
	store := tempBaselineStore(t)
	_ = store.Set("production", "snap-001")
	_ = store.Set("production", "snap-002")
	label, err := store.Get("production")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if label != "snap-002" {
		t.Errorf("expected snap-002, got %q", label)
	}
}

func TestBaselineStore_Remove(t *testing.T) {
	store := tempBaselineStore(t)
	_ = store.Set("staging", "snap-010")
	if err := store.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := store.Get("staging")
	if err == nil {
		t.Fatal("expected error after removal, got nil")
	}
}

func TestBaselineStore_RemoveMissing(t *testing.T) {
	store := tempBaselineStore(t)
	err := store.Remove("nonexistent")
	if err == nil {
		t.Fatal("expected error removing nonexistent entry, got nil")
	}
}

func TestBaselineStore_List(t *testing.T) {
	store := tempBaselineStore(t)
	_ = store.Set("production", "snap-001")
	_ = store.Set("staging", "snap-010")
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestBaselineStore_ListEmpty(t *testing.T) {
	store := tempBaselineStore(t)
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
