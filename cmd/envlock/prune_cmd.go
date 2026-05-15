package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envlock/internal/snapshot"
)

func runPrune(args []string, store *snapshot.Store, out io.Writer) error {
	fs := flag.NewFlagSet("prune", flag.ContinueOnError)
	fs.SetOutput(out)

	keepLast := fs.Int("keep-last", 0, "retain the N most recent snapshots")
	olderThanDays := fs.Int("older-than-days", 0, "remove snapshots older than N days")
	dryRun := fs.Bool("dry-run", false, "report what would be pruned without deleting")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *keepLast == 0 && *olderThanDays == 0 {
		return fmt.Errorf("prune: at least one of --keep-last or --older-than-days is required")
	}

	opts := snapshot.PruneOptions{
		KeepLast: *keepLast,
		DryRun:   *dryRun,
	}
	if *olderThanDays > 0 {
		opts.OlderThan = time.Now().AddDate(0, 0, -*olderThanDays)
	}

	result, err := snapshot.Prune(store, opts)
	if err != nil {
		return err
	}

	prefix := ""
	if *dryRun {
		prefix = "[dry-run] "
	}

	if len(result.Removed) == 0 {
		fmt.Fprintf(out, "%snothing to prune\n", prefix)
		return nil
	}

	for _, l := range result.Removed {
		fmt.Fprintf(out, "%spruned: %s\n", prefix, l)
	}
	fmt.Fprintf(out, "%s%d snapshot(s) pruned, %d retained\n", prefix, len(result.Removed), len(result.Retained))
	return nil
}

func init() {
	_ = os.Stderr // ensure os is used
}
