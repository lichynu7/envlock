package snapshot

import "strings"

// FilterOptions controls which environment variables are included or excluded.
type FilterOptions struct {
	// PrefixAllowlist, if non-empty, only includes vars with one of these prefixes.
	PrefixAllowlist []string
	// PrefixBlocklist excludes vars with any of these prefixes.
	PrefixBlocklist []string
	// ExactBlocklist excludes vars with these exact names.
	ExactBlocklist []string
}

// DefaultSensitiveBlocklist returns a blocklist of commonly sensitive variable names.
func DefaultSensitiveBlocklist() []string {
	return []string{
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
		"GITHUB_TOKEN",
		"NPM_TOKEN",
		"DOCKER_PASSWORD",
		"DATABASE_URL",
		"SECRET_KEY",
		"PRIVATE_KEY",
	}
}

// Apply filters a map of environment variables according to the given options.
// It returns a new map containing only the variables that pass all filters.
func (f FilterOptions) Apply(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))

	exactBlock := make(map[string]struct{}, len(f.ExactBlocklist))
	for _, name := range f.ExactBlocklist {
		exactBlock[name] = struct{}{}
	}

OUTER:
	for k, v := range env {
		// Exact blocklist check
		if _, blocked := exactBlock[k]; blocked {
			continue
		}

		// Prefix blocklist check
		for _, prefix := range f.PrefixBlocklist {
			if strings.HasPrefix(k, prefix) {
				continue OUTER
			}
		}

		// Prefix allowlist check (if set, variable must match at least one)
		if len(f.PrefixAllowlist) > 0 {
			allowed := false
			for _, prefix := range f.PrefixAllowlist {
				if strings.HasPrefix(k, prefix) {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}

		result[k] = v
	}

	return result
}
