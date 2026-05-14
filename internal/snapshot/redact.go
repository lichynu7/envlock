package snapshot

import "strings"

// RedactOptions controls how sensitive values are redacted in snapshots.
type RedactOptions struct {
	// Blocklist is a list of exact key names to redact.
	Blocklist []string
	// PrefixBlocklist is a list of key prefixes to redact.
	PrefixBlocklist []string
	// Placeholder is the string used to replace redacted values.
	// Defaults to "[REDACTED]" if empty.
	Placeholder string
}

// DefaultRedactOptions returns a RedactOptions using the DefaultSensitiveBlocklist
// and common sensitive prefixes.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Blocklist: DefaultSensitiveBlocklist,
		PrefixBlocklist: []string{
			"SECRET_",
			"PRIVATE_",
			"TOKEN_",
			"CREDENTIAL_",
		},
		Placeholder: "[REDACTED]",
	}
}

// Redact returns a new Snapshot with sensitive variable values replaced by the
// placeholder string. The original Snapshot is not modified.
func Redact(s Snapshot, opts RedactOptions) Snapshot {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	blocklist := make(map[string]struct{}, len(opts.Blocklist))
	for _, k := range opts.Blocklist {
		blocklist[k] = struct{}{}
	}

	redacted := make(map[string]string, len(s.Vars))
	for k, v := range s.Vars {
		if _, blocked := blocklist[k]; blocked {
			redacted[k] = placeholder
			continue
		}
		upper := strings.ToUpper(k)
		matched := false
		for _, prefix := range opts.PrefixBlocklist {
			if strings.HasPrefix(upper, strings.ToUpper(prefix)) {
				matched = true
				break
			}
		}
		if matched {
			redacted[k] = placeholder
		} else {
			redacted[k] = v
		}
	}

	return Snapshot{
		Label:     s.Label,
		Timestamp: s.Timestamp,
		Vars:      redacted,
	}
}
