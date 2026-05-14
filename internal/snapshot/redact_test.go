package snapshot

import (
	"testing"
)

func baseRedactEnv() map[string]string {
	return map[string]string{
		"HOME":              "/home/user",
		"PATH":              "/usr/bin:/bin",
		"AWS_SECRET_ACCESS_KEY": "supersecret",
		"DATABASE_PASSWORD": "dbpass",
		"SECRET_KEY":        "mykey",
		"TOKEN_GITHUB":      "ghtoken",
		"PRIVATE_KEY":       "privkey",
		"APP_NAME":          "envlock",
	}
}

func TestRedact_BlocklistExactMatch(t *testing.T) {
	s := Snapshot{Label: "test", Vars: baseRedactEnv()}
	opts := RedactOptions{
		Blocklist:   []string{"DATABASE_PASSWORD"},
		Placeholder: "[REDACTED]",
	}
	out := Redact(s, opts)
	if out.Vars["DATABASE_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DATABASE_PASSWORD to be redacted, got %q", out.Vars["DATABASE_PASSWORD"])
	}
	if out.Vars["APP_NAME"] != "envlock" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out.Vars["APP_NAME"])
	}
}

func TestRedact_PrefixBlocklist(t *testing.T) {
	s := Snapshot{Label: "test", Vars: baseRedactEnv()}
	opts := RedactOptions{
		PrefixBlocklist: []string{"SECRET_", "TOKEN_"},
		Placeholder:     "[REDACTED]",
	}
	out := Redact(s, opts)
	if out.Vars["SECRET_KEY"] != "[REDACTED]" {
		t.Errorf("expected SECRET_KEY to be redacted")
	}
	if out.Vars["TOKEN_GITHUB"] != "[REDACTED]" {
		t.Errorf("expected TOKEN_GITHUB to be redacted")
	}
	if out.Vars["HOME"] != "/home/user" {
		t.Errorf("expected HOME to be unchanged")
	}
}

func TestRedact_DefaultOptions(t *testing.T) {
	s := Snapshot{Label: "test", Vars: baseRedactEnv()}
	opts := DefaultRedactOptions()
	out := Redact(s, opts)
	if out.Vars["AWS_SECRET_ACCESS_KEY"] != "[REDACTED]" {
		t.Errorf("expected AWS_SECRET_ACCESS_KEY to be redacted, got %q", out.Vars["AWS_SECRET_ACCESS_KEY"])
	}
	if out.Vars["TOKEN_GITHUB"] != "[REDACTED]" {
		t.Errorf("expected TOKEN_GITHUB to be redacted")
	}
	if out.Vars["APP_NAME"] != "envlock" {
		t.Errorf("expected APP_NAME to be unchanged")
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	vars := baseRedactEnv()
	s := Snapshot{Label: "test", Vars: vars}
	opts := DefaultRedactOptions()
	_ = Redact(s, opts)
	if s.Vars["AWS_SECRET_ACCESS_KEY"] != "supersecret" {
		t.Errorf("original snapshot was mutated")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	s := Snapshot{Label: "test", Vars: map[string]string{"DATABASE_PASSWORD": "secret"}}
	opts := RedactOptions{
		Blocklist:   []string{"DATABASE_PASSWORD"},
		Placeholder: "***",
	}
	out := Redact(s, opts)
	if out.Vars["DATABASE_PASSWORD"] != "***" {
		t.Errorf("expected custom placeholder, got %q", out.Vars["DATABASE_PASSWORD"])
	}
}

func TestRedact_EmptyPlaceholderUsesDefault(t *testing.T) {
	s := Snapshot{Label: "test", Vars: map[string]string{"DATABASE_PASSWORD": "secret"}}
	opts := RedactOptions{
		Blocklist:   []string{"DATABASE_PASSWORD"},
		Placeholder: "",
	}
	out := Redact(s, opts)
	if out.Vars["DATABASE_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected default placeholder, got %q", out.Vars["DATABASE_PASSWORD"])
	}
}
