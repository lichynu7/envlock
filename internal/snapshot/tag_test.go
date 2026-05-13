package snapshot_test

import (
	"os"
	"testing"

	"github.com/user/envlock/internal/snapshot"
)

func tempTagStore(t *testing.T) *snapshot.TagStore {
	t.Helper()
	dir, err := os.MkdirTemp("", "envlock-tags-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return snapshot.NewTagStore(dir)
}

func TestTagStore_AddAndLabels(t *testing.T) {
	ts := tempTagStore(t)
	if err := ts.Add("release", "snap-1"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := ts.Add("release", "snap-2"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	labels, err := ts.Labels("release")
	if err != nil {
		t.Fatalf("Labels: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestTagStore_AddIdempotent(t *testing.T) {
	ts := tempTagStore(t)
	_ = ts.Add("ci", "snap-1")
	_ = ts.Add("ci", "snap-1")
	labels, _ := ts.Labels("ci")
	if len(labels) != 1 {
		t.Errorf("expected 1 label after duplicate add, got %d", len(labels))
	}
}

func TestTagStore_Remove(t *testing.T) {
	ts := tempTagStore(t)
	_ = ts.Add("staging", "snap-a")
	_ = ts.Add("staging", "snap-b")
	if err := ts.Remove("staging", "snap-a"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	labels, _ := ts.Labels("staging")
	if len(labels) != 1 || labels[0] != "snap-b" {
		t.Errorf("expected [snap-b], got %v", labels)
	}
}

func TestTagStore_RemoveLastLabelDeletesTag(t *testing.T) {
	ts := tempTagStore(t)
	_ = ts.Add("tmp", "snap-x")
	_ = ts.Remove("tmp", "snap-x")
	tags, _ := ts.Tags()
	for _, tag := range tags {
		if tag == "tmp" {
			t.Error("expected tag 'tmp' to be removed")
		}
	}
}

func TestTagStore_Tags(t *testing.T) {
	ts := tempTagStore(t)
	_ = ts.Add("beta", "snap-1")
	_ = ts.Add("alpha", "snap-2")
	tags, err := ts.Tags()
	if err != nil {
		t.Fatalf("Tags: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "alpha" || tags[1] != "beta" {
		t.Errorf("expected sorted tags [alpha beta], got %v", tags)
	}
}

func TestTagStore_EmptyLabels(t *testing.T) {
	ts := tempTagStore(t)
	labels, err := ts.Labels("nonexistent")
	if err != nil {
		t.Fatalf("Labels: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty labels, got %v", labels)
	}
}
