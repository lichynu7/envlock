package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// HistoryEntry represents a single snapshot entry in the history list.
type HistoryEntry struct {
	Label     string
	Timestamp time.Time
	VarCount  int
}

// HistoryOptions controls how history is retrieved and displayed.
type HistoryOptions struct {
	// Limit caps the number of entries returned. 0 means no limit.
	Limit int
	// Prefix filters entries to only those whose label starts with the given string.
	Prefix string
	// SortDesc sorts entries from newest to oldest when true (default: oldest first).
	SortDesc bool
}

// History returns a list of snapshot history entries from the given store,
// applying any options provided. Entries are sorted by timestamp.
func History(store *Store, opts HistoryOptions) ([]HistoryEntry, error) {
	labels, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("history: listing snapshots: %w", err)
	}

	entries := make([]HistoryEntry, 0, len(labels))
	for _, label := range labels {
		// Apply prefix filter if set.
		if opts.Prefix != "" && len(label) < len(opts.Prefix) {
			continue
		}
		if opts.Prefix != "" && label[:len(opts.Prefix)] != opts.Prefix {
			continue
		}

		snap, err := store.Load(label)
		if err != nil {
			// Skip snapshots that can't be loaded rather than failing entirely.
			continue
		}

		entries = append(entries, HistoryEntry{
			Label:     snap.Label,
			Timestamp: snap.Timestamp,
			VarCount:  len(snap.Vars),
		})
	}

	// Sort by timestamp.
	sort.Slice(entries, func(i, j int) bool {
		if opts.SortDesc {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		}
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	// Apply limit.
	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[:opts.Limit]
	}

	return entries, nil
}

// FormatHistory renders a slice of HistoryEntry values as a human-readable
// table string suitable for terminal output.
func FormatHistory(entries []HistoryEntry) string {
	if len(entries) == 0 {
		return "No snapshots found.\n"
	}

	// Header
	result := fmt.Sprintf("%-30s  %-25s  %s\n", "LABEL", "TIMESTAMP", "VARS")
	result += fmt.Sprintf("%-30s  %-25s  %s\n",
		"------------------------------",
		"-------------------------",
		"----")

	for _, e := range entries {
		label := e.Label
		if len(label) > 30 {
			label = label[:27] + "..."
		}
		result += fmt.Sprintf("%-30s  %-25s  %d\n",
			label,
			e.Timestamp.Format(time.RFC3339),
			e.VarCount,
		)
	}

	return result
}
