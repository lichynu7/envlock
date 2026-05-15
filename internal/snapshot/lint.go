package snapshot

import (
	"fmt"
	"strings"
)

// LintRule describes a single lint check applied to a snapshot.
type LintRule struct {
	Name    string
	Message string
}

// LintResult holds the outcome of linting a snapshot.
type LintResult struct {
	Label    string
	Warnings []LintRule
}

// IsClean returns true when no warnings were found.
func (r LintResult) IsClean() bool {
	return len(r.Warnings) == 0
}

// LintOptions controls which rules are applied during linting.
type LintOptions struct {
	// WarnOnEmpty flags variables whose value is an empty string.
	WarnOnEmpty bool
	// WarnOnLowercaseKeys flags variables that are not fully uppercase.
	WarnOnLowercaseKeys bool
	// WarnOnLongValues flags values that exceed MaxValueLength bytes.
	WarnOnLongValues bool
	MaxValueLength   int
}

// DefaultLintOptions returns a sensible default configuration.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		WarnOnEmpty:         true,
		WarnOnLowercaseKeys: true,
		WarnOnLongValues:    true,
		MaxValueLength:      512,
	}
}

// Lint analyses a snapshot and returns a LintResult containing any warnings.
func Lint(s Snapshot, opts LintOptions) LintResult {
	result := LintResult{Label: s.Label}

	for key, val := range s.Vars {
		if opts.WarnOnEmpty && val == "" {
			result.Warnings = append(result.Warnings, LintRule{
				Name:    "empty-value",
				Message: fmt.Sprintf("%s has an empty value", key),
			})
		}

		if opts.WarnOnLowercaseKeys && key != strings.ToUpper(key) {
			result.Warnings = append(result.Warnings, LintRule{
				Name:    "lowercase-key",
				Message: fmt.Sprintf("%s contains lowercase characters", key),
			})
		}

		if opts.WarnOnLongValues && len(val) > opts.MaxValueLength {
			result.Warnings = append(result.Warnings, LintRule{
				Name:    "long-value",
				Message: fmt.Sprintf("%s value exceeds %d bytes (%d)", key, opts.MaxValueLength, len(val)),
			})
		}
	}

	return result
}
