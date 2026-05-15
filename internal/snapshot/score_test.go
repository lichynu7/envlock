package snapshot

import (
	"strings"
	"testing"
)

func makeScoreSnapshot(vars map[string]string) Snapshot {
	return Snapshot{Label: "test", Vars: vars}
}

func TestScore_PerfectScore(t *testing.T) {
	s := makeScoreSnapshot(map[string]string{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
	})
	r := Score(s, DefaultScoreOptions())
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
	if len(r.Penalties) != 0 {
		t.Errorf("expected no penalties, got %v", r.Penalties)
	}
}

func TestScore_PenalisesEmptyValue(t *testing.T) {
	s := makeScoreSnapshot(map[string]string{"APP_ENV": ""})
	opts := DefaultScoreOptions()
	r := Score(s, opts)
	if r.Score != 100-opts.PenaltyEmptyValue {
		t.Errorf("expected %d, got %d", 100-opts.PenaltyEmptyValue, r.Score)
	}
	if len(r.Penalties) != 1 {
		t.Errorf("expected 1 penalty, got %d", len(r.Penalties))
	}
}

func TestScore_PenalisesLowercaseKey(t *testing.T) {
	s := makeScoreSnapshot(map[string]string{"app_env": "dev"})
	opts := DefaultScoreOptions()
	r := Score(s, opts)
	if r.Score != 100-opts.PenaltyLowercaseKey {
		t.Errorf("expected %d, got %d", 100-opts.PenaltyLowercaseKey, r.Score)
	}
}

func TestScore_PenalisesLongValue(t *testing.T) {
	long := strings.Repeat("x", 300)
	s := makeScoreSnapshot(map[string]string{"BIG_VAR": long})
	opts := DefaultScoreOptions()
	r := Score(s, opts)
	if r.Score != 100-opts.PenaltyLongValue {
		t.Errorf("expected %d, got %d", 100-opts.PenaltyLongValue, r.Score)
	}
}

func TestScore_PenalisesSensitiveKey(t *testing.T) {
	s := makeScoreSnapshot(map[string]string{"AWS_SECRET_ACCESS_KEY": "abc123"})
	opts := DefaultScoreOptions()
	r := Score(s, opts)
	if r.Score != 100-opts.PenaltySensitiveKey {
		t.Errorf("expected %d, got %d", 100-opts.PenaltySensitiveKey, r.Score)
	}
}

func TestScore_NeverBelowZero(t *testing.T) {
	long := strings.Repeat("z", 300)
	vars := map[string]string{
		"aws_secret_access_key": "",
		"db_password":           long,
	}
	s := makeScoreSnapshot(vars)
	opts := DefaultScoreOptions()
	opts.PenaltyEmptyValue = 50
	opts.PenaltyLowercaseKey = 50
	opts.PenaltyLongValue = 50
	opts.PenaltySensitiveKey = 50
	r := Score(s, opts)
	if r.Score < 0 {
		t.Errorf("score should not be negative, got %d", r.Score)
	}
}

func TestScore_LabelPreserved(t *testing.T) {
	s := makeScoreSnapshot(map[string]string{})
	s.Label = "my-snapshot"
	r := Score(s, DefaultScoreOptions())
	if r.Label != "my-snapshot" {
		t.Errorf("expected label my-snapshot, got %s", r.Label)
	}
}
