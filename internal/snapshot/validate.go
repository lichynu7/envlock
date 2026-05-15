package snapshot

import (
	"fmt"
	"strings"
)

// ValidationResult holds the outcome of validating a single environment variable.
type ValidationResult struct {
	Key     string
	Message string
}

// ValidationOptions controls which checks are applied during validation.
type ValidationOptions struct {
	// RequireUppercase enforces that all keys are uppercase.
	RequireUppercase bool
	// ForbidEmpty rejects variables with empty values.
	ForbidEmpty bool
	// ForbidPrefix rejects keys that start with any of the given prefixes.
	ForbidPrefix []string
	// RequirePrefix requires that every key starts with one of the given prefixes.
	RequirePrefix []string
}

// DefaultValidateOptions returns a sensible default configuration.
func DefaultValidateOptions() ValidationOptions {
	return ValidationOptions{
		RequireUppercase: true,
		ForbidEmpty:      true,
	}
}

// Validate checks each variable in the snapshot against the supplied options
// and returns a slice of ValidationResult for every violation found.
func Validate(s Snapshot, opts ValidationOptions) []ValidationResult {
	var results []ValidationResult

	for k, v := range s.Vars {
		if opts.RequireUppercase && k != strings.ToUpper(k) {
			results = append(results, ValidationResult{
				Key:     k,
				Message: "key must be uppercase",
			})
		}

		if opts.ForbidEmpty && v == "" {
			results = append(results, ValidationResult{
				Key:     k,
				Message: "value must not be empty",
			})
		}

		for _, prefix := range opts.ForbidPrefix {
			if strings.HasPrefix(k, prefix) {
				results = append(results, ValidationResult{
					Key:     k,
					Message: fmt.Sprintf("key must not start with forbidden prefix %q", prefix),
				})
				break
			}
		}

		if len(opts.RequirePrefix) > 0 {
			matched := false
			for _, prefix := range opts.RequirePrefix {
				if strings.HasPrefix(k, prefix) {
					matched = true
					break
				}
			}
			if !matched {
				results = append(results, ValidationResult{
					Key:     k,
					Message: fmt.Sprintf("key must start with one of the required prefixes %v", opts.RequirePrefix),
				})
			}
		}
	}

	return results
}
