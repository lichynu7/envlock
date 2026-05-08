package snapshot

import "fmt"

// DiffEntry represents a single changed, added, or removed environment variable.
type DiffEntry struct {
	Key      string
	OldValue string
	NewValue string
	Status   DiffStatus
}

// DiffStatus indicates the type of change for a variable.
type DiffStatus string

const (
	StatusAdded   DiffStatus = "added"
	StatusRemoved DiffStatus = "removed"
	StatusChanged DiffStatus = "changed"
)

// String returns a human-readable representation of a DiffEntry.
func (d DiffEntry) String() string {
	switch d.Status {
	case StatusAdded:
		return fmt.Sprintf("+ %s=%s", d.Key, d.NewValue)
	case StatusRemoved:
		return fmt.Sprintf("- %s=%s", d.Key, d.OldValue)
	case StatusChanged:
		return fmt.Sprintf("~ %s: %q -> %q", d.Key, d.OldValue, d.NewValue)
	default:
		return fmt.Sprintf("? %s", d.Key)
	}
}

// Diff compares two snapshots and returns a slice of DiffEntry describing
// the variables that were added, removed, or changed between them.
func Diff(base, target *Snapshot) []DiffEntry {
	var entries []DiffEntry

	baseMap := make(map[string]string, len(base.Vars))
	for _, v := range base.Vars {
		baseMap[v.Key] = v.Value
	}

	targetMap := make(map[string]string, len(target.Vars))
	for _, v := range target.Vars {
		targetMap[v.Key] = v.Value
	}

	// Detect added and changed variables.
	for k, newVal := range targetMap {
		oldVal, exists := baseMap[k]
		if !exists {
			entries = append(entries, DiffEntry{Key: k, NewValue: newVal, Status: StatusAdded})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{Key: k, OldValue: oldVal, NewValue: newVal, Status: StatusChanged})
		}
	}

	// Detect removed variables.
	for k, oldVal := range baseMap {
		if _, exists := targetMap[k]; !exists {
			entries = append(entries, DiffEntry{Key: k, OldValue: oldVal, Status: StatusRemoved})
		}
	}

	return entries
}
