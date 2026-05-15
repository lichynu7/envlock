package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// SummaryOptions controls what is included in a snapshot summary.
type SummaryOptions struct {
	TopN        int  // number of longest values to highlight
	ShowPrefixes bool // group keys by common prefix
}

// DefaultSummaryOptions returns sensible defaults for SummaryOptions.
func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		TopN:        5,
		ShowPrefixes: true,
	}
}

// Summary holds aggregated statistics about a snapshot.
type Summary struct {
	Label       string
	TotalVars   int
	EmptyVars   int
	PrefixGroups map[string]int // prefix -> count
	LongestVars []string       // keys of the TopN longest values
}

// Summarize produces a Summary from a Snapshot using the given options.
func Summarize(s Snapshot, opts SummaryOptions) Summary {
	sum := Summary{
		Label:        s.Label,
		TotalVars:    len(s.Vars),
		PrefixGroups: make(map[string]int),
	}

	type kv struct {
		key string
		len int
	}
	var lengths []kv

	for k, v := range s.Vars {
		if v == "" {
			sum.EmptyVars++
		}
		lengths = append(lengths, kv{k, len(v)})

		if opts.ShowPrefixes {
			if idx := strings.Index(k, "_"); idx > 0 {
				prefix := k[:idx]
				sum.PrefixGroups[prefix]++
			} else {
				sum.PrefixGroups["(none)"]++
			}
		}
	}

	sort.Slice(lengths, func(i, j int) bool {
		return lengths[i].len > lengths[j].len
	})

	n := opts.TopN
	if n > len(lengths) {
		n = len(lengths)
	}
	for i := 0; i < n; i++ {
		sum.LongestVars = append(sum.LongestVars, lengths[i].key)
	}

	return sum
}

// FormatSummary renders a Summary as a human-readable string.
func FormatSummary(sum Summary) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot : %s\n", sum.Label)
	fmt.Fprintf(&sb, "Total    : %d vars\n", sum.TotalVars)
	fmt.Fprintf(&sb, "Empty    : %d vars\n", sum.EmptyVars)

	if len(sum.PrefixGroups) > 0 {
		prefixes := make([]string, 0, len(sum.PrefixGroups))
		for p := range sum.PrefixGroups {
			prefixes = append(prefixes, p)
		}
		sort.Strings(prefixes)
		sb.WriteString("Prefixes :\n")
		for _, p := range prefixes {
			fmt.Fprintf(&sb, "  %-20s %d\n", p, sum.PrefixGroups[p])
		}
	}

	if len(sum.LongestVars) > 0 {
		sb.WriteString("Longest  :\n")
		for _, k := range sum.LongestVars {
			fmt.Fprintf(&sb, "  %s\n", k)
		}
	}

	return sb.String()
}
