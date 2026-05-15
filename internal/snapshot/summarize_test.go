package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeSummarizeSnapshot(vars map[string]string) Snapshot {
	return Snapshot{
		Label:     "test",
		Timestamp: time.Now(),
		Vars:      vars,
	}
}

func TestSummarize_TotalAndEmpty(t *testing.T) {
	s := makeSummarizeSnapshot(map[string]string{
		"HOME":  "/root",
		"EMPTY": "",
		"PATH":  "/usr/bin",
	})
	sum := Summarize(s, DefaultSummaryOptions())
	if sum.TotalVars != 3 {
		t.Errorf("expected 3 total vars, got %d", sum.TotalVars)
	}
	if sum.EmptyVars != 1 {
		t.Errorf("expected 1 empty var, got %d", sum.EmptyVars)
	}
}

func TestSummarize_PrefixGroups(t *testing.T) {
	s := makeSummarizeSnapshot(map[string]string{
		"AWS_REGION":     "us-east-1",
		"AWS_ACCESS_KEY": "AKID",
		"HOME":           "/root",
	})
	opts := DefaultSummaryOptions()
	sum := Summarize(s, opts)

	if sum.PrefixGroups["AWS"] != 2 {
		t.Errorf("expected AWS prefix count 2, got %d", sum.PrefixGroups["AWS"])
	}
	if sum.PrefixGroups["(none)"] != 1 {
		t.Errorf("expected (none) prefix count 1, got %d", sum.PrefixGroups["(none)"])
	}
}

func TestSummarize_LongestVars(t *testing.T) {
	s := makeSummarizeSnapshot(map[string]string{
		"A": "short",
		"B": strings.Repeat("x", 200),
		"C": strings.Repeat("y", 100),
		"D": "tiny",
	})
	opts := DefaultSummaryOptions()
	opts.TopN = 2
	sum := Summarize(s, opts)

	if len(sum.LongestVars) != 2 {
		t.Fatalf("expected 2 longest vars, got %d", len(sum.LongestVars))
	}
	if sum.LongestVars[0] != "B" {
		t.Errorf("expected B as longest, got %s", sum.LongestVars[0])
	}
	if sum.LongestVars[1] != "C" {
		t.Errorf("expected C as second longest, got %s", sum.LongestVars[1])
	}
}

func TestSummarize_NoPrefixes(t *testing.T) {
	s := makeSummarizeSnapshot(map[string]string{
		"HOME": "/root",
		"USER": "alice",
	})
	opts := DefaultSummaryOptions()
	opts.ShowPrefixes = false
	sum := Summarize(s, opts)

	if len(sum.PrefixGroups) != 0 {
		t.Errorf("expected no prefix groups, got %d", len(sum.PrefixGroups))
	}
}

func TestFormatSummary_ContainsLabel(t *testing.T) {
	s := makeSummarizeSnapshot(map[string]string{
		"FOO_BAR": "baz",
	})
	sum := Summarize(s, DefaultSummaryOptions())
	out := FormatSummary(sum)

	if !strings.Contains(out, "test") {
		t.Errorf("expected label 'test' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Total") {
		t.Errorf("expected 'Total' in output, got:\n%s", out)
	}
}
