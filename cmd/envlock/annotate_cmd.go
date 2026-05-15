package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/envlock/internal/snapshot"
)

// runAnnotate handles the `envlock annotate` subcommand.
// Usage:
//
//	envlock annotate set   <label> <note>
//	envlock annotate get   <label>
//	envlock annotate remove <label>
//	envlock annotate list
func runAnnotate(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("annotate", flag.ContinueOnError)
	fs.SetOutput(out)
	dir := fs.String("store", snapshot.DefaultStore(), "annotation store directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		return fmt.Errorf("annotate: subcommand required (set|get|remove|list)")
	}

	store, err := snapshot.NewAnnotationStore(*dir)
	if err != nil {
		return err
	}

	switch remaining[0] {
	case "set":
		if len(remaining) < 3 {
			return fmt.Errorf("annotate set: usage: annotate set <label> <note>")
		}
		if err := store.Set(remaining[1], remaining[2]); err != nil {
			return err
		}
		fmt.Fprintf(out, "annotation set for %q\n", remaining[1])

	case "get":
		if len(remaining) < 2 {
			return fmt.Errorf("annotate get: usage: annotate get <label>")
		}
		a, err := store.Get(remaining[1])
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "label:      %s\nnote:       %s\ncreated_at: %s\nupdated_at: %s\n",
			a.Label, a.Note, a.CreatedAt.Format("2006-01-02T15:04:05Z"), a.UpdatedAt.Format("2006-01-02T15:04:05Z"))

	case "remove":
		if len(remaining) < 2 {
			return fmt.Errorf("annotate remove: usage: annotate remove <label>")
		}
		if err := store.Remove(remaining[1]); err != nil {
			return err
		}
		fmt.Fprintf(out, "annotation removed for %q\n", remaining[1])

	case "list":
		annotations, err := store.List()
		if err != nil {
			return err
		}
		if len(annotations) == 0 {
			fmt.Fprintln(out, "no annotations")
			return nil
		}
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "LABEL\tNOTE\tUPDATED")
		for _, a := range annotations {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", a.Label, a.Note, a.UpdatedAt.Format("2006-01-02"))
		}
		tw.Flush()

	default:
		return fmt.Errorf("annotate: unknown subcommand %q", remaining[0])
	}
	return nil
}

func init() {
	_ = os.Stderr // ensure os import used
}
