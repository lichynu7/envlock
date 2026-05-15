package snapshot

import (
	"testing"
)

func makeValidateSnapshot(vars map[string]string) Snapshot {
	return Snapshot{
		Label: "validate-test",
		Vars:  vars,
	}
}

func TestValidate_NoViolations(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin",
	})
	results := Validate(s, DefaultValidateOptions())
	if len(results) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(results), results)
	}
}

func TestValidate_LowercaseKey(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"home": "/home/user",
	})
	opts := DefaultValidateOptions()
	opts.ForbidEmpty = false
	results := Validate(s, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Key != "home" {
		t.Errorf("expected violation for key 'home', got %q", results[0].Key)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"MY_VAR": "",
	})
	opts := DefaultValidateOptions()
	opts.RequireUppercase = false
	results := Validate(s, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Message != "value must not be empty" {
		t.Errorf("unexpected message: %q", results[0].Message)
	}
}

func TestValidate_ForbidPrefix(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"SECRET_TOKEN": "abc123",
		"APP_NAME":     "envlock",
	})
	opts := ValidationOptions{
		ForbidPrefix: []string{"SECRET_"},
	}
	results := Validate(s, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(results), results)
	}
	if results[0].Key != "SECRET_TOKEN" {
		t.Errorf("expected violation for SECRET_TOKEN, got %q", results[0].Key)
	}
}

func TestValidate_RequirePrefix(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"APP_HOST": "localhost",
		"DATABASE": "postgres",
	})
	opts := ValidationOptions{
		RequirePrefix: []string{"APP_"},
	}
	results := Validate(s, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(results), results)
	}
	if results[0].Key != "DATABASE" {
		t.Errorf("expected violation for DATABASE, got %q", results[0].Key)
	}
}

func TestValidate_MultipleViolationsSameKey(t *testing.T) {
	s := makeValidateSnapshot(map[string]string{
		"bad_key": "",
	})
	results := Validate(s, DefaultValidateOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 violations (lowercase + empty), got %d: %+v", len(results), results)
	}
}
