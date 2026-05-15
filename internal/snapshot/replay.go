package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// ReplayEntry represents a single step in a replay sequence.
type ReplayEntry struct {
	Step      int
	Label     string
	Timestamp time.Time
	VarCount  int
	Diff      *DiffResult
}

// ReplayOptions controls how a replay is constructed.
type ReplayOptions struct {
	// IncludeDiff attaches a diff from the previous snapshot to each entry.
	IncludeDiff bool
}

// DefaultReplayOptions returns sensible defaults.
func DefaultReplayOptions() ReplayOptions {
	return ReplayOptions{
		IncludeDiff: true,
	}
}

// Replay builds an ordered sequence of ReplayEntry values from a store,
// optionally computing the diff between consecutive snapshots.
func Replay(store *Store, opts ReplayOptions) ([]ReplayEntry, error) {
	labels, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("replay: list snapshots: %w", err)
	}
	if len(labels) == 0 {
		return nil, nil
	}

	snaps := make([]Snapshot, 0, len(labels))
	for _, label := range labels {
		s, err := store.Load(label)
		if err != nil {
			return nil, fmt.Errorf("replay: load %q: %w", label, err)
		}
		snaps = append(snaps, s)
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].Timestamp.Before(snaps[j].Timestamp)
	})

	entries := make([]ReplayEntry, 0, len(snaps))
	for i, s := range snaps {
		entry := ReplayEntry{
			Step:      i + 1,
			Label:     s.Label,
			Timestamp: s.Timestamp,
			VarCount:  len(s.Vars),
		}
		if opts.IncludeDiff && i > 0 {
			dr := Diff(snaps[i-1], s)
			entry.Diff = &dr
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// FormatReplay returns a human-readable string summarising a replay sequence.
func FormatReplay(entries []ReplayEntry) string {
	if len(entries) == 0 {
		return "no snapshots to replay\n"
	}
	out := ""
	for _, e := range entries {
		out += fmt.Sprintf("step %d: %s (%s) — %d vars",
			e.Step, e.Label, e.Timestamp.Format(time.RFC3339), e.VarCount)
		if e.Diff != nil {
			out += fmt.Sprintf(" [+%d -%d ~%d]",
				len(e.Diff.Added), len(e.Diff.Removed), len(e.Diff.Changed))
		}
		out += "\n"
	}
	return out
}
