package snapshot

import (
	"time"
)

// WatchEvent represents a change detected during environment watching.
type WatchEvent struct {
	At      time.Time
	Diff    []DiffEntry
	Baseline *Snapshot
	Current  *Snapshot
}

// WatchOptions configures the Watch function.
type WatchOptions struct {
	// Interval is how often to poll for changes.
	Interval time.Duration
	// FilterOpts is applied to each captured snapshot before diffing.
	FilterOpts *FilterOptions
	// Label prefix used for internally named snapshots.
	Label string
}

// DefaultWatchOptions returns sensible defaults for watching.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval: 5 * time.Second,
		Label:    "watch-baseline",
	}
}

// Watch captures a baseline snapshot then emits WatchEvents on the returned
// channel whenever the environment changes. It stops when the done channel is
// closed. The caller is responsible for draining the events channel.
func Watch(opts WatchOptions, done <-chan struct{}) (<-chan WatchEvent, error) {
	baseline, err := Capture(opts.Label)
	if err != nil {
		return nil, err
	}

	if opts.FilterOpts != nil {
		baseline.Vars = Filter(baseline.Vars, *opts.FilterOpts)
	}

	events := make(chan WatchEvent, 8)

	go func() {
		defer close(events)
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		current := baseline
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				next, err := Capture(opts.Label)
				if err != nil {
					continue
				}
				if opts.FilterOpts != nil {
					next.Vars = Filter(next.Vars, *opts.FilterOpts)
				}
				entries := Diff(current, next)
				if len(entries) > 0 {
					events <- WatchEvent{
						At:       t,
						Diff:     entries,
						Baseline: current,
						Current:  next,
					}
					current = next
				}
			}
		}
	}()

	return events, nil
}
