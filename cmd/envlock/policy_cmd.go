package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

// runPolicy handles the "policy" subcommand.
// Usage:
//
//	envlock policy --label <label> --policy <file> [--store <dir>]
func runPolicy(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("policy", flag.ContinueOnError)
	fs.SetOutput(out)

	label := fs.String("label", "", "snapshot label to validate")
	policyFile := fs.String("policy", "", "path to policy JSON file")
	storeDir := fs.String("store", defaultStoreDir(), "snapshot store directory")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *label == "" {
		return fmt.Errorf("--label is required")
	}
	if *policyFile == "" {
		return fmt.Errorf("--policy is required")
	}

	store := snapshot.DefaultStore(*storeDir)
	snap, err := store.Load(*label)
	if err != nil {
		return fmt.Errorf("load snapshot %q: %w", *label, err)
	}

	pol, err := snapshot.LoadPolicy(*policyFile)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}

	violations := snapshot.EnforcePolicy(snap, pol)
	if len(violations) == 0 {
		fmt.Fprintf(out, "OK: snapshot %q satisfies policy %q\n", *label, pol.Name)
		return nil
	}

	fmt.Fprintf(out, "FAIL: %d violation(s) for policy %q on snapshot %q:\n",
		len(violations), pol.Name, *label)
	for _, v := range violations {
		fmt.Fprintf(out, "  [%s] %s\n", v.Key, v.Message)
	}
	os.Exit(1)
	return nil
}
