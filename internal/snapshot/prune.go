package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// PruneOptions controls which snapshots are removed.
type PruneOptions struct {
	// KeepLast retains the N most recent snapshots. 0 means no limit.
	KeepLast int
	// OlderThan removes snapshots created before this time. Zero value is ignored.
	OlderThan time.Time
	// DryRun reports what would be pruned without deleting.
	DryRun bool
}

// PruneResult summarises the outcome of a prune operation.
type PruneResult struct {
	Removed []string
	Retained []string
}

// Prune removes snapshots from store according to opts.
// Labels are sorted by timestamp descending so the newest are retained first.
func Prune(store *Store, opts PruneOptions) (PruneResult, error) {
	labels, err := store.List()
	if err != nil {
		return PruneResult{}, fmt.Errorf("prune: list snapshots: %w", err)
	}

	type entry struct {
		label string
		ts    time.Time
	}

	entries := make([]entry, 0, len(labels))
	for _, l := range labels {
		snap, err := store.Load(l)
		if err != nil {
			return PruneResult{}, fmt.Errorf("prune: load %q: %w", l, err)
		}
		entries = append(entries, entry{label: l, ts: snap.Timestamp})
	}

	// Sort newest first.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ts.After(entries[j].ts)
	})

	var result PruneResult
	for i, e := range entries {
		shouldRemove := false

		if opts.KeepLast > 0 && i >= opts.KeepLast {
			shouldRemove = true
		}
		if !opts.OlderThan.IsZero() && e.ts.Before(opts.OlderThan) {
			shouldRemove = true
		}

		if shouldRemove {
			result.Removed = append(result.Removed, e.label)
			if !opts.DryRun {
				if err := store.Delete(e.label); err != nil {
					return result, fmt.Errorf("prune: delete %q: %w", e.label, err)
				}
			}
		} else {
			result.Retained = append(result.Retained, e.label)
		}
	}
	return result, nil
}
