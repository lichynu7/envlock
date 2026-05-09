package main

import (
	"fmt"
	"os"

	"github.com/user/envlock/internal/snapshot"
)

// runExport handles the `envlock export <label> [--format=<fmt>]` command.
func runExport(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envlock export <label> [--format json|csv|dotenv]")
	}

	label := args[0]
	format := snapshot.FormatDotEnv

	for _, arg := range args[1:] {
		switch arg {
		case "--format=json", "-f=json":
			format = snapshot.FormatJSON
		case "--format=csv", "-f=csv":
			format = snapshot.FormatCSV
		case "--format=dotenv", "-f=dotenv":
			format = snapshot.FormatDotEnv
		default:
			return fmt.Errorf("unknown flag: %q", arg)
		}
	}

	store := snapshot.DefaultStore()
	s, err := store.Load(label)
	if err != nil {
		return fmt.Errorf("could not load snapshot %q: %w", label, err)
	}

	if err := snapshot.Export(s, format, os.Stdout); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	return nil
}
