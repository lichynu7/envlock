package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envlock/internal/snapshot"
)

func runCompare(args []string, store *snapshot.Store, out io.Writer) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	verbose := fs.Bool("verbose", false, "show full diff in addition to summary")
	fs.SetOutput(out)

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("compare requires two snapshot labels: compare <labelA> <labelB>")
	}

	labelA := remaining[0]
	labelB := remaining[1]

	snapshotA, err := store.Load(labelA)
	if err != nil {
		return fmt.Errorf("could not load snapshot '%s': %w", labelA, err)
	}

	snapshotB, err := store.Load(labelB)
	if err != nil {
		return fmt.Errorf("could not load snapshot '%s': %w", labelB, err)
	}

	result := snapshot.Compare(snapshotA, snapshotB)
	fmt.Fprintln(out, result.Summary())

	if *verbose {
		entries := snapshot.SortedDiff(snapshotA, snapshotB)
		if len(entries) > 0 {
			fmt.Fprintln(out)
			fmt.Fprint(out, snapshot.RenderDiff(entries))
		}
	}

	if result.Added > 0 || result.Removed > 0 || result.Changed > 0 {
		os.Exit(1)
	}

	return nil
}
