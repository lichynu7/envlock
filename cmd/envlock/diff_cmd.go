package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlock/internal/snapshot"
)

func runDiff(args []string) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	colorFlag := fs.Bool("color", true, "enable colored output")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("usage: envlock diff <label-a> <label-b>")
	}

	labelA := remaining[0]
	labelB := remaining[1]

	store, err := snapshot.DefaultStore()
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}

	snapshotA, err := store.Load(labelA)
	if err != nil {
		return fmt.Errorf("could not load snapshot %q: %w", labelA, err)
	}

	snapshotB, err := store.Load(labelB)
	if err != nil {
		return fmt.Errorf("could not load snapshot %q: %w", labelB, err)
	}

	result := snapshot.Diff(snapshotA, snapshotB)
	output := snapshot.RenderDiff(result, *colorFlag)

	fmt.Print(output)
	return nil
}
