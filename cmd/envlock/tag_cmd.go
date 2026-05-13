package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlock/internal/snapshot"
)

func runTag(args []string, out io.Writer, storeDir string) error {
	fs := flag.NewFlagSet("tag", flag.ContinueOnError)
	fs.SetOutput(out)
	remove := fs.Bool("remove", false, "remove the tag from the snapshot")
	list := fs.Bool("list", false, "list all tags or labels for a tag")

	if err := fs.Parse(args); err != nil {
		return err
	}

	ts := snapshot.NewTagStore(storeDir)

	if *list {
		remaining := fs.Args()
		if len(remaining) == 0 {
			tags, err := ts.Tags()
			if err != nil {
				return fmt.Errorf("listing tags: %w", err)
			}
			if len(tags) == 0 {
				fmt.Fprintln(out, "no tags found")
				return nil
			}
			fmt.Fprintln(out, strings.Join(tags, "\n"))
			return nil
		}
		tag := remaining[0]
		labels, err := ts.Labels(tag)
		if err != nil {
			return fmt.Errorf("listing labels for tag %q: %w", tag, err)
		}
		if len(labels) == 0 {
			fmt.Fprintf(out, "no snapshots tagged %q\n", tag)
			return nil
		}
		fmt.Fprintln(out, strings.Join(labels, "\n"))
		return nil
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("usage: envlock tag [--remove] <tag> <label>")
	}
	tag, label := remaining[0], remaining[1]

	if *remove {
		if err := ts.Remove(tag, label); err != nil {
			return fmt.Errorf("removing tag: %w", err)
		}
		fmt.Fprintf(out, "removed tag %q from snapshot %q\n", tag, label)
		return nil
	}

	if err := ts.Add(tag, label); err != nil {
		return fmt.Errorf("adding tag: %w", err)
	}
	fmt.Fprintf(out, "tagged snapshot %q with %q\n", label, tag)
	return nil
}
