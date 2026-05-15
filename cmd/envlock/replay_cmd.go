package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

func runReplay(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("replay", flag.ContinueOnError)
	fs.SetOutput(out)

	noDiff := fs.Bool("no-diff", false, "omit per-step diff summary")
	storeDir := fs.String("store", snapshot.DefaultStore(), "path to snapshot store")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	store, err := snapshot.NewStore(*storeDir)
	if err != nil {
		fmt.Fprintf(out, "error: open store: %v\n", err)
		return 1
	}

	opts := snapshot.DefaultReplayOptions()
	if *noDiff {
		opts.IncludeDiff = false
	}

	entries, err := snapshot.Replay(store, opts)
	if err != nil {
		fmt.Fprintf(out, "error: %v\n", err)
		return 1
	}

	fmt.Fprint(out, snapshot.FormatReplay(entries))
	return 0
}

func init() {
	_ = os.Stderr // ensure os is used
}
