package snapshot

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// HistoryEntry represents a single entry in the snapshot history.
type HistoryEntry struct {
	Label     string
	Timestamp time.Time
	VarCount  int
}

// History returns a list of HistoryEntry values for all snapshots in the store,
// sorted by timestamp descending (most recent first).
func History(store *Store) ([]HistoryEntry, error) {
	labels, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("history: listing snapshots: %w", err)
	}

	entries := make([]HistoryEntry, 0, len(labels))
	for _, label := range labels {
		snap, err := store.Load(label)
		if err != nil {
			return nil, fmt.Errorf("history: loading snapshot %q: %w", label, err)
		}
		entries = append(entries, HistoryEntry{
			Label:     snap.Label,
			Timestamp: snap.Timestamp,
			VarCount:  len(snap.Vars),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	return entries, nil
}

// FormatHistory renders a list of HistoryEntry values as a human-readable table string.
func FormatHistory(entries []HistoryEntry) string {
	if len(entries) == 0 {
		return "No snapshots found.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %-25s  %s\n", "LABEL", "TIMESTAMP", "VARS"))
	sb.WriteString(strings.Repeat("-", 65) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-30s  %-25s  %d\n",
			truncate(e.Label, 30),
			e.Timestamp.Format(time.RFC3339),
			e.VarCount,
		))
	}
	return sb.String()
}
