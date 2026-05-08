package snapshot

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// OutputFormat controls how diff results are rendered.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

// RenderDiff writes a human-readable diff to w.
func RenderDiff(w io.Writer, entries []DiffEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	// Sort entries: removed, added, changed — then alphabetically.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			order := map[DiffKind]int{Removed: 0, Added: 1, Changed: 2}
			return order[entries[i].Kind] < order[entries[j].Kind]
		}
		return entries[i].Key < entries[j].Key
	})

	for _, e := range entries {
		switch e.Kind {
		case Added:
			fmt.Fprintf(tw, "+ %s\t= %s\n", e.Key, truncate(e.NewValue, 80))
		case Removed:
			fmt.Fprintf(tw, "- %s\t= %s\n", e.Key, truncate(e.OldValue, 80))
		case Changed:
			fmt.Fprintf(tw, "~ %s\t%s -> %s\n", e.Key, truncate(e.OldValue, 40), truncate(e.NewValue, 40))
		}
	}
	return nil
}

// RenderSnapshot writes a snapshot's variables as key=value lines to w.
func RenderSnapshot(w io.Writer, s *Snapshot) error {
	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintf(tw, "Label:\t%s\n", s.Label)
	fmt.Fprintf(tw, "Captured:\t%s\n", s.CapturedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(tw, "Variables:\t%d\n", len(s.Vars))
	fmt.Fprintln(tw, strings.Repeat("-", 40))
	for _, k := range keys {
		fmt.Fprintf(tw, "  %s\t= %s\n", k, truncate(s.Vars[k], 80))
	}
	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
