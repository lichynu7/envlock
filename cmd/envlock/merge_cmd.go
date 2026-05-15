package main

import (
	"flag"
	"fmt"
	"io"
	"sort"

	"github.com/user/envlock/internal/snapshot"
)

func runMerge(args []string, store *snapshot.Store, out io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(out)

	var label string
	var strategy string

	fs.StringVar(&label, "label", "", "label for the merged snapshot (default: auto-generated)")
	fs.StringVar(&strategy, "strategy", "base", "conflict resolution strategy: base or overlay")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envlock merge [--label NAME] [--strategy base|overlay] <base-label> <overlay-label>")
	}

	baseLabel := fs.Arg(0)
	overlayLabel := fs.Arg(1)

	base, err := store.Load(baseLabel)
	if err != nil {
		return fmt.Errorf("load base snapshot %q: %w", baseLabel, err)
	}

	overlay, err := store.Load(overlayLabel)
	if err != nil {
		return fmt.Errorf("load overlay snapshot %q: %w", overlayLabel, err)
	}

	var mergeStrategy snapshot.MergeStrategy
	switch strategy {
	case "overlay":
		mergeStrategy = snapshot.MergeStrategyOverlay
	case "base", "":
		mergeStrategy = snapshot.MergeStrategyBase
	default:
		return fmt.Errorf("unknown strategy %q: must be 'base' or 'overlay'", strategy)
	}

	result := snapshot.Merge(base, overlay, snapshot.MergeOptions{
		Strategy: mergeStrategy,
		Label:    label,
	})

	if err := store.Save(result.Snapshot); err != nil {
		return fmt.Errorf("save merged snapshot: %w", err)
	}

	fmt.Fprintf(out, "Merged snapshot saved as %q\n", result.Snapshot.Label)
	fmt.Fprintf(out, "  vars: %d | conflicts: %d | base-only: %d | overlay-only: %d\n",
		len(result.Snapshot.Vars), len(result.Conflicts),
		len(result.BaseOnly), len(result.OverlayOnly))

	if len(result.Conflicts) > 0 {
		sort.Strings(result.Conflicts)
		fmt.Fprintf(out, "  conflict keys (resolved via %s):\n", strategy)
		for _, k := range result.Conflicts {
			fmt.Fprintf(out, "    - %s\n", k)
		}
	}

	return nil
}
