package snapshot

import (
	"testing"
	"time"
)

// TestWatch_NoChanges verifies that no events are emitted when the environment
// does not change between polls.
func TestWatch_NoChanges(t *testing.T) {
	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		Label:    "test-watch",
	}

	done := make(chan struct{})
	events, err := Watch(opts, done)
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	// Let it tick a couple of times then stop.
	time.Sleep(60 * time.Millisecond)
	close(done)

	var received []WatchEvent
	for e := range events {
		received = append(received, e)
	}

	// The environment should not have changed, so we expect no events.
	if len(received) != 0 {
		t.Errorf("expected 0 events, got %d", len(received))
	}
}

// TestWatch_EventHasFields verifies that a WatchEvent is well-formed.
func TestWatch_EventHasFields(t *testing.T) {
	baseline, err := Capture("evt-base")
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	current, err := Capture("evt-cur")
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}

	// Manufacture a synthetic event.
	evt := WatchEvent{
		At:       time.Now(),
		Diff:     Diff(baseline, current),
		Baseline: baseline,
		Current:  current,
	}

	if evt.At.IsZero() {
		t.Error("expected non-zero At timestamp")
	}
	if evt.Baseline == nil {
		t.Error("expected non-nil Baseline")
	}
	if evt.Current == nil {
		t.Error("expected non-nil Current")
	}
}

// TestDefaultWatchOptions verifies sane defaults.
func TestDefaultWatchOptions(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Interval <= 0 {
		t.Errorf("expected positive Interval, got %v", opts.Interval)
	}
	if opts.Label == "" {
		t.Error("expected non-empty Label")
	}
}

// TestWatch_DoneClosesChannel ensures the events channel is closed after done.
func TestWatch_DoneClosesChannel(t *testing.T) {
	opts := WatchOptions{
		Interval: 10 * time.Millisecond,
		Label:    "test-done",
	}

	done := make(chan struct{})
	events, err := Watch(opts, done)
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	close(done)

	select {
	case _, ok := <-events:
		if ok {
			// drain remaining
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("events channel was not closed after done")
	}
}
