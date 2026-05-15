package main

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envlock/internal/snapshot"
)

func runBaseline(args []string, out io.Writer, storeDir string) error {
	fs := flag.NewFlagSet("baseline", flag.ContinueOnError)
	fs.SetOutput(out)

	var remove bool
	fs.BoolVar(&remove, "remove", false, "remove the baseline for the given environment")

	if err := fs.Parse(args); err != nil {
		return err
	}

	store, err := snapshot.NewBaselineStore(storeDir)
	if err != nil {
		return fmt.Errorf("baseline: %w", err)
	}

	parsed := fs.Args()

	// List mode: no positional arguments
	if len(parsed) == 0 {
		entries, err := store.List()
		if err != nil {
			return fmt.Errorf("baseline list: %w", err)
		}
		if len(entries) == 0 {
			fmt.Fprintln(out, "No baselines set.")
			return nil
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Environment < entries[j].Environment
		})
		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ENVIRONMENT\tLABEL")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\n", e.Environment, e.Label)
		}
		return w.Flush()
	}

	if len(parsed) < 2 {
		return fmt.Errorf("baseline: usage: baseline [--remove] <environment> <label>")
	}

	environment := parsed[0]
	label := parsed[1]

	if remove {
		if err := store.Remove(environment); err != nil {
			return fmt.Errorf("baseline remove: %w", err)
		}
		fmt.Fprintf(out, "Baseline removed for environment %q.\n", environment)
		return nil
	}

	if err := store.Set(environment, label); err != nil {
		return fmt.Errorf("baseline set: %w", err)
	}
	fmt.Fprintf(out, "Baseline for %q set to %q.\n", environment, label)
	return nil
}
