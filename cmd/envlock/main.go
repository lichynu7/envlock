// main.go is the entry point for the envlock CLI tool.
// It wires together the available subcommands and dispatches
// based on the first argument provided by the user.
package main

import (
	"fmt"
	"os"

	"github.com/envlock/internal/snapshot"
)

const usage = `envlock — snapshot and diff environment variables

Usage:
  envlock <command> [options]

Commands:
  capture   Capture the current environment and save it with a label
  diff      Diff two saved environment snapshots
  export    Export a snapshot to JSON, CSV, or .env format
  list      List all saved snapshots

Run 'envlock <command> -help' for command-specific options.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	store := snapshot.DefaultStore()

	subcmd := os.Args[1]
	args := os.Args[2:]

	var err error

	switch subcmd {
	case "capture":
		err = runCapture(store, args)
	case "diff":
		err = runDiff(store, args)
	case "export":
		err = runExport(store, args)
	case "list":
		err = runList(store, args)
	case "-help", "--help", "-h", "help":
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "envlock: unknown command %q\n\n", subcmd)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "envlock %s: %v\n", subcmd, err)
		os.Exit(1)
	}
}
