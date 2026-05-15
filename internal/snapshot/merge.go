package snapshot

import "fmt"

// MergeStrategy controls how conflicts are resolved when merging two snapshots.
type MergeStrategy int

const (
	// MergeStrategyBase keeps the value from the base snapshot on conflict.
	MergeStrategyBase MergeStrategy = iota
	// MergeStrategyOverlay keeps the value from the overlay snapshot on conflict.
	MergeStrategyOverlay
)

// MergeOptions configures the Merge operation.
type MergeOptions struct {
	Strategy MergeStrategy
	Label    string
}

// MergeResult holds the merged snapshot and metadata about the merge.
type MergeResult struct {
	Snapshot   Snapshot
	Conflicts  []string // keys that had conflicting values
	BaseOnly   []string // keys only in base
	OverlayOnly []string // keys only in overlay
}

// Merge combines two snapshots into one according to the given options.
// The base snapshot provides the foundation; the overlay is merged on top.
func Merge(base, overlay Snapshot, opts MergeOptions) MergeResult {
	merged := make(map[string]string)
	var conflicts, baseOnly, overlayOnly []string

	for k, v := range base.Vars {
		merged[k] = v
	}

	for k, ov := range overlay.Vars {
		bv, exists := merged[k]
		if !exists {
			merged[k] = ov
			overlayOnly = append(overlayOnly, k)
			continue
		}
		if bv != ov {
			conflicts = append(conflicts, k)
			if opts.Strategy == MergeStrategyOverlay {
				merged[k] = ov
			}
		}
	}

	for k := range base.Vars {
		if _, ok := overlay.Vars[k]; !ok {
			baseOnly = append(baseOnly, k)
		}
	}

	label := opts.Label
	if label == "" {
		label = fmt.Sprintf("merge(%s,%s)", base.Label, overlay.Label)
	}

	return MergeResult{
		Snapshot:    Snapshot{Label: label, Vars: merged},
		Conflicts:   conflicts,
		BaseOnly:    baseOnly,
		OverlayOnly: overlayOnly,
	}
}
