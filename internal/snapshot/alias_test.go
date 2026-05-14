package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func tempAliasStore(t *testing.T) *AliasStore {
	t.Helper()
	dir := t.TempDir()
	as, err := NewAliasStore(filepath.Join(dir, "aliases.json"))
	if err != nil {
		t.Fatalf("NewAliasStore: %v", err)
	}
	return as
}

func TestAliasStore_SetAndResolve(t *testing.T) {
	as := tempAliasStore(t)
	if err := as.Set("prod", "snapshot-2024-prod"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	label, err := as.Resolve("prod")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if label != "snapshot-2024-prod" {
		t.Errorf("expected snapshot-2024-prod, got %s", label)
	}
}

func TestAliasStore_ResolveUnknown(t *testing.T) {
	as := tempAliasStore(t)
	_, err := as.Resolve("nonexistent")
	if err == nil {
		t.Error("expected error for unknown alias")
	}
}

func TestAliasStore_Remove(t *testing.T) {
	as := tempAliasStore(t)
	_ = as.Set("staging", "snapshot-staging-01")
	if err := as.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := as.Resolve("staging")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestAliasStore_RemoveUnknown(t *testing.T) {
	as := tempAliasStore(t)
	if err := as.Remove("ghost"); err == nil {
		t.Error("expected error removing unknown alias")
	}
}

func TestAliasStore_List(t *testing.T) {
	as := tempAliasStore(t)
	_ = as.Set("a", "label-a")
	_ = as.Set("b", "label-b")
	list := as.List()
	if len(list) != 2 {
		t.Errorf("expected 2 entries, got %d", len(list))
	}
	if list["a"] != "label-a" || list["b"] != "label-b" {
		t.Errorf("unexpected list contents: %v", list)
	}
}

func TestAliasStore_Persistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")

	as1, _ := NewAliasStore(path)
	_ = as1.Set("ci", "snapshot-ci-42")

	as2, err := NewAliasStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	label, err := as2.Resolve("ci")
	if err != nil {
		t.Fatalf("Resolve after reload: %v", err)
	}
	if label != "snapshot-ci-42" {
		t.Errorf("expected snapshot-ci-42, got %s", label)
	}
}

func TestAliasStore_SetEmptyAlias(t *testing.T) {
	as := tempAliasStore(t)
	if err := as.Set("", "some-label"); err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestAliasStore_MissingFileIsOK(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "aliases.json")
	_, err := NewAliasStore(path)
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
	_ = os.Remove(path)
}
