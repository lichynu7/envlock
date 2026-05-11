package snapshot

import (
	"strings"
	"testing"
)

func makeCompareSnapshot(label string, vars map[string]string) *Snapshot {
	return &Snapshot{Label: label, Vars: vars}
}

func TestCompare_NoDiff(t *testing.T) {
	a := makeCompareSnapshot("base", map[string]string{"FOO": "1", "BAR": "2"})
	b := makeCompareSnapshot("head", map[string]string{"FOO": "1", "BAR": "2"})

	result := Compare(a, b)

	if result.Added != 0 || result.Removed != 0 || result.Changed != 0 {
		t.Errorf("expected no diff, got %+v", result)
	}
	if result.Total != 2 {
		t.Errorf("expected Total=2, got %d", result.Total)
	}
}

func TestCompare_CountsAdded(t *testing.T) {
	a := makeCompareSnapshot("base", map[string]string{"FOO": "1"})
	b := makeCompareSnapshot("head", map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"})

	result := Compare(a, b)

	if result.Added != 2 {
		t.Errorf("expected Added=2, got %d", result.Added)
	}
	if result.Removed != 0 || result.Changed != 0 {
		t.Errorf("unexpected diffs: %+v", result)
	}
}

func TestCompare_CountsRemoved(t *testing.T) {
	a := makeCompareSnapshot("base", map[string]string{"FOO": "1", "BAR": "2"})
	b := makeCompareSnapshot("head", map[string]string{"FOO": "1"})

	result := Compare(a, b)

	if result.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", result.Removed)
	}
}

func TestCompare_CountsChanged(t *testing.T) {
	a := makeCompareSnapshot("base", map[string]string{"FOO": "old"})
	b := makeCompareSnapshot("head", map[string]string{"FOO": "new"})

	result := Compare(a, b)

	if result.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", result.Changed)
	}
}

func TestCompareResult_Summary_NoDiff(t *testing.T) {
	r := CompareResult{LabelA: "base", LabelB: "head"}
	s := r.Summary()
	if !strings.Contains(s, "No differences") {
		t.Errorf("expected no-diff message, got: %s", s)
	}
}

func TestCompareResult_Summary_WithDiff(t *testing.T) {
	r := CompareResult{LabelA: "base", LabelB: "head", Added: 1, Removed: 2, Changed: 3, Total: 10}
	s := r.Summary()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "-2") || !strings.Contains(s, "~3") {
		t.Errorf("unexpected summary format: %s", s)
	}
}

func TestSortedDiff_OrderedByKey(t *testing.T) {
	a := makeCompareSnapshot("base", map[string]string{"ZZZ": "1", "AAA": "1", "MMM": "1"})
	b := makeCompareSnapshot("head", map[string]string{"ZZZ": "2", "AAA": "2", "MMM": "2"})

	entries := SortedDiff(a, b)

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "AAA" || entries[1].Key != "MMM" || entries[2].Key != "ZZZ" {
		t.Errorf("entries not sorted: %v", entries)
	}
}
