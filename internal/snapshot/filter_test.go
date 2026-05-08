package snapshot

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOME":                 "/home/user",
		"PATH":                 "/usr/bin:/bin",
		"AWS_REGION":           "us-east-1",
		"AWS_SECRET_ACCESS_KEY": "supersecret",
		"GITHUB_TOKEN":         "ghp_abc123",
		"MY_APP_PORT":          "8080",
		"MY_APP_ENV":           "production",
	}
}

func TestFilter_ExactBlocklist(t *testing.T) {
	f := FilterOptions{
		ExactBlocklist: []string{"GITHUB_TOKEN", "AWS_SECRET_ACCESS_KEY"},
	}
	result := f.Apply(baseEnv())

	if _, ok := result["GITHUB_TOKEN"]; ok {
		t.Error("expected GITHUB_TOKEN to be blocked")
	}
	if _, ok := result["AWS_SECRET_ACCESS_KEY"]; ok {
		t.Error("expected AWS_SECRET_ACCESS_KEY to be blocked")
	}
	if _, ok := result["HOME"]; !ok {
		t.Error("expected HOME to be present")
	}
}

func TestFilter_PrefixBlocklist(t *testing.T) {
	f := FilterOptions{
		PrefixBlocklist: []string{"AWS_"},
	}
	result := f.Apply(baseEnv())

	for k := range result {
		if len(k) >= 4 && k[:4] == "AWS_" {
			t.Errorf("expected %q to be blocked by prefix", k)
		}
	}
	if _, ok := result["HOME"]; !ok {
		t.Error("expected HOME to be present")
	}
}

func TestFilter_PrefixAllowlist(t *testing.T) {
	f := FilterOptions{
		PrefixAllowlist: []string{"MY_APP_"},
	}
	result := f.Apply(baseEnv())

	if len(result) != 2 {
		t.Errorf("expected 2 vars, got %d", len(result))
	}
	if _, ok := result["MY_APP_PORT"]; !ok {
		t.Error("expected MY_APP_PORT to be present")
	}
	if _, ok := result["MY_APP_ENV"]; !ok {
		t.Error("expected MY_APP_ENV to be present")
	}
}

func TestFilter_NoOptions(t *testing.T) {
	f := FilterOptions{}
	env := baseEnv()
	result := f.Apply(env)

	if len(result) != len(env) {
		t.Errorf("expected all %d vars, got %d", len(env), len(result))
	}
}

func TestFilter_DefaultSensitiveBlocklist(t *testing.T) {
	blocklist := DefaultSensitiveBlocklist()
	if len(blocklist) == 0 {
		t.Error("expected non-empty default sensitive blocklist")
	}

	f := FilterOptions{ExactBlocklist: blocklist}
	result := f.Apply(baseEnv())

	if _, ok := result["AWS_SECRET_ACCESS_KEY"]; ok {
		t.Error("expected AWS_SECRET_ACCESS_KEY to be filtered by default blocklist")
	}
}
