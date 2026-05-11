package snapshot

import (
	"fmt"
	"sort"
)

// CompareResult holds a summary of differences between two snapshots.
type CompareResult struct {
	LabelA  string
	LabelB  string
	Added   int
	Removed int
	Changed int
	Total   int
}

// Summary returns a human-readable one-line summary of the comparison.
func (c CompareResult) Summary() string {
	if c.Added == 0 && c.Removed == 0 && c.Changed == 0 {
		return fmt.Sprintf("No differences found between '%s' and '%s'.", c.LabelA, c.LabelB)
	}
	return fmt.Sprintf(
		"'%s' vs '%s': +%d added, -%d removed, ~%d changed (%d total vars)",
		c.LabelA, c.LabelB, c.Added, c.Removed, c.Changed, c.Total,
	)
}

// Compare runs a diff between two snapshots and returns a CompareResult.
func Compare(a, b *Snapshot) CompareResult {
	entries := Diff(a, b)

	result := CompareResult{
		LabelA: a.Label,
		LabelB: b.Label,
	}

	allKeys := make(map[string]struct{})
	for k := range a.Vars {
		allKeys[k] = struct{}{}
	}
	for k := range b.Vars {
		allKeys[k] = struct{}{}
	}
	result.Total = len(allKeys)

	for _, e := range entries {
		switch e.Kind {
		case DiffAdded:
			result.Added++
		case DiffRemoved:
			result.Removed++
		case DiffChanged:
			result.Changed++
		}
	}

	return result
}

// SortedDiff returns diff entries sorted by key name for deterministic output.
func SortedDiff(a, b *Snapshot) []DiffEntry {
	entries := Diff(a, b)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
