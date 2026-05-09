package snapshot

import (
	"strings"
	"testing"
	"time"
)

func TestRenderDiff_Added(t *testing.T) {
	entries := []DiffEntry{
		{Key: "NEW_VAR", Kind: Added, NewValue: "hello"},
	}
	var buf strings.Builder
	if err := RenderDiff(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ NEW_VAR") {
		t.Errorf("expected added marker, got: %s", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected value in output, got: %s", out)
	}
}

func TestRenderDiff_Removed(t *testing.T) {
	entries := []DiffEntry{
		{Key: "OLD_VAR", Kind: Removed, OldValue: "bye"},
	}
	var buf strings.Builder
	if err := RenderDiff(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- OLD_VAR") {
		t.Errorf("expected removed marker, got: %s", out)
	}
}

func TestRenderDiff_Changed(t *testing.T) {
	entries := []DiffEntry{
		{Key: "PATH", Kind: Changed, OldValue: "/usr/bin", NewValue: "/usr/local/bin"},
	}
	var buf strings.Builder
	if err := RenderDiff(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ PATH") {
		t.Errorf("expected changed marker, got: %s", out)
	}
	if !strings.Contains(out, "->") {
		t.Errorf("expected arrow separator, got: %s", out)
	}
}

func TestRenderDiff_Empty(t *testing.T) {
	var buf strings.Builder
	if err := RenderDiff(&buf, []DiffEntry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestRenderDiff_MultipleEntries(t *testing.T) {
	entries := []DiffEntry{
		{Key: "ADDED_VAR", Kind: Added, NewValue: "new"},
		{Key: "REMOVED_VAR", Kind: Removed, OldValue: "old"},
		{Key: "CHANGED_VAR", Kind: Changed, OldValue: "before", NewValue: "after"},
	}
	var buf strings.Builder
	if err := RenderDiff(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ ADDED_VAR") {
		t.Errorf("expected added entry in output, got: %s", out)
	}
	if !strings.Contains(out, "- REMOVED_VAR") {
		t.Errorf("expected removed entry in output, got: %s", out)
	}
	if !strings.Contains(out, "~ CHANGED_VAR") {
		t.Errorf("expected changed entry in output, got: %s", out)
	}
}

func TestRenderSnapshot_ContainsLabel(t *testing.T) {
	s := &Snapshot{
		Label:      "ci-run",
		CapturedAt: time.Now(),
		Vars:       map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
	var buf strings.Builder
	if err := RenderSnapshot(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ci-run") {
		t.Errorf("expected label in output, got: %s", out)
	}
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "bar") {
		t.Errorf("expected variable FOO=bar in output, got: %s", out)
	}
}

func TestTruncate(t *testing.T) {
	long := strings.Repeat("x", 100)
	result := truncate(long, 80)
	if len(result) > 83 {
		t.Errorf("expected truncated string, got length %d", len(result))
	}
	if !strings.HasSuffix(result, "...") {
		t.Errorf("expected ellipsis suffix, got: %s", result)
	}

	short := "hello"
	if truncate(short, 80) != short {
		t.Errorf("expected unchanged short string")
	}
}
