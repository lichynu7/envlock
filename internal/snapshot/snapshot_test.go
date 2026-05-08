package snapshot_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/envlock/internal/snapshot"
)

func TestCapture_ContainsKnownVar(t *testing.T) {
	const key = "ENVLOCK_TEST_VAR"
	const val = "hello_envlock"
	t.Setenv(key, val)

	s := snapshot.Capture("test")

	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if s.Vars[key] != val {
		t.Errorf("expected %q=%q, got %q", key, val, s.Vars[key])
	}
}

func TestCapture_Label(t *testing.T) {
	s := snapshot.Capture("my-label")
	if s.Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", s.Label)
	}
}

func TestCapture_TimestampSet(t *testing.T) {
	s := snapshot.Capture("")
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	os.Setenv("ENVLOCK_ROUND_TRIP", "42")
	defer os.Unsetenv("ENVLOCK_ROUND_TRIP")

	orig := snapshot.Capture("roundtrip")
	data, err := orig.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	restored, err := snapshot.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if restored.Label != orig.Label {
		t.Errorf("label mismatch: got %q, want %q", restored.Label, orig.Label)
	}
	if restored.Vars["ENVLOCK_ROUND_TRIP"] != "42" {
		t.Errorf("var mismatch after round-trip")
	}
}

func TestMarshal_ValidJSON(t *testing.T) {
	s := snapshot.Capture("json-test")
	data, err := s.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}

func TestSortedKeys_Order(t *testing.T) {
	s := &snapshot.Snapshot{
		Vars: map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	keys := s.SortedKeys()
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("key[%d]: got %q, want %q", i, k, expected[i])
		}
	}
}
