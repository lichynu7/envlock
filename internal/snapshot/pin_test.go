package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPinStore(t *testing.T) *PinStore {
	t.Helper()
	dir := t.TempDir()
	return NewPinStore(filepath.Join(dir, "pins.json"))
}

func TestPin_AddAndList(t *testing.T) {
	store := tempPinStore(t)
	if err := store.Pin("API_KEY", "abc123", "main api key"); err != nil {
		t.Fatalf("Pin failed: %v", err)
	}
	pins, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	if pins[0].Key != "API_KEY" || pins[0].Value != "abc123" {
		t.Errorf("unexpected pin: %+v", pins[0])
	}
}

func TestPin_Idempotent(t *testing.T) {
	store := tempPinStore(t)
	_ = store.Pin("FOO", "v1", "")
	_ = store.Pin("FOO", "v2", "updated")
	pins, _ := store.List()
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	if pins[0].Value != "v2" {
		t.Errorf("expected updated value v2, got %s", pins[0].Value)
	}
}

func TestPin_Unpin(t *testing.T) {
	store := tempPinStore(t)
	_ = store.Pin("REMOVE_ME", "val", "")
	if err := store.Unpin("REMOVE_ME"); err != nil {
		t.Fatalf("Unpin failed: %v", err)
	}
	pins, _ := store.List()
	if len(pins) != 0 {
		t.Errorf("expected 0 pins after unpin, got %d", len(pins))
	}
}

func TestPin_UnpinMissing(t *testing.T) {
	store := tempPinStore(t)
	if err := store.Unpin("NONEXISTENT"); err == nil {
		t.Error("expected error when unpinning missing key")
	}
}

func TestPin_ListEmpty(t *testing.T) {
	store := tempPinStore(t)
	pins, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected empty list, got %d", len(pins))
	}
}

func TestPin_CheckDrift_NoDrift(t *testing.T) {
	store := tempPinStore(t)
	_ = store.Pin("DB_HOST", "localhost", "")
	env := map[string]string{"DB_HOST": "localhost", "OTHER": "val"}
	drifted, err := store.CheckDrift(env)
	if err != nil {
		t.Fatalf("CheckDrift failed: %v", err)
	}
	if len(drifted) != 0 {
		t.Errorf("expected no drift, got: %v", drifted)
	}
}

func TestPin_CheckDrift_Detected(t *testing.T) {
	store := tempPinStore(t)
	_ = store.Pin("DB_HOST", "localhost", "")
	_ = store.Pin("PORT", "5432", "")
	env := map[string]string{"DB_HOST": "remotehost"}  // PORT is missing too
	drifted, err := store.CheckDrift(env)
	if err != nil {
		t.Fatalf("CheckDrift failed: %v", err)
	}
	if len(drifted) != 2 {
		t.Errorf("expected 2 drifted keys, got %d: %v", len(drifted), drifted)
	}
}

func TestPin_PersistsAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	s1 := NewPinStore(path)
	_ = s1.Pin("PERSIST_ME", "hello", "")
	s2 := NewPinStore(path)
	pins, err := s2.List()
	if err != nil {
		t.Fatalf("List on second store failed: %v", err)
	}
	if len(pins) != 1 || pins[0].Key != "PERSIST_ME" {
		t.Errorf("pin did not persist across store instances")
	}
	_ = os.Remove(path)
}
