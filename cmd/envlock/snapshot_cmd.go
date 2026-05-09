package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/envlock/internal/snapshot"
)

// runSnapshot handles the 'snapshot' subcommand, which captures the current
// environment variables and saves them under a given label.
//
// Usage:
//
//	envlock snapshot -label <name> [options]
//
// Flags:
//
//	-label      Required. Name to associate with this snapshot.
//	-no-filter  Disable the default sensitive variable blocklist.
func runSnapshot(args []string, store *snapshot.Store, out io.Writer) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	fs.SetOutput(out)

	label := fs.String("label", "", "Label to assign to this snapshot (required)")
	noFilter := fs.Bool("no-filter", false, "Disable filtering of sensitive environment variables")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *label == "" {
		fs.Usage()
		return fmt.Errorf("flag -label is required")
	}

	// Build filter options based on flags.
	var filterOpts []snapshot.FilterOption
	if !*noFilter {
		filterOpts = append(filterOpts, snapshot.WithBlocklist(snapshot.DefaultSensitiveBlocklist))
	}

	// Capture the current environment.
	env := os.Environ()
	filtered := snapshot.Filter(env, filterOpts...)

	snap := snapshot.Capture(*label, filtered)

	if err := store.Save(snap); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	fmt.Fprintf(out, "Snapshot %q saved (%d variables)\n", snap.Label, len(snap.Vars))
	return nil
}
