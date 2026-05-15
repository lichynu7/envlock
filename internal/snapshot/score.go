package snapshot

import "fmt"

// ScoreResult holds the computed health score for a snapshot.
type ScoreResult struct {
	Label      string
	Score      int // 0–100
	Penalties  []string
}

// ScoreOptions controls how the score is calculated.
type ScoreOptions struct {
	PenaltyEmptyValue    int
	PenaltyLowercaseKey  int
	PenaltyLongValue     int
	PenaltySensitiveKey  int
	LongValueThreshold   int
}

// DefaultScoreOptions returns sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		PenaltyEmptyValue:   5,
		PenaltyLowercaseKey: 3,
		PenaltyLongValue:    4,
		PenaltySensitiveKey: 10,
		LongValueThreshold:  256,
	}
}

// Score evaluates the quality of a snapshot's environment variables
// and returns a ScoreResult with a 0–100 health score.
func Score(s Snapshot, opts ScoreOptions) ScoreResult {
	result := ScoreResult{Label: s.Label, Score: 100}

	sensitive := make(map[string]bool)
	for _, k := range DefaultSensitiveBlocklist() {
		sensitive[k] = true
	}

	for k, v := range s.Vars {
		if v == "" {
			result.Penalties = append(result.Penalties,
				fmt.Sprintf("%s: empty value (-%d)", k, opts.PenaltyEmptyValue))
			result.Score -= opts.PenaltyEmptyValue
		}
		if len(v) > opts.LongValueThreshold {
			result.Penalties = append(result.Penalties,
				fmt.Sprintf("%s: value too long (-%d)", k, opts.PenaltyLongValue))
			result.Score -= opts.PenaltyLongValue
		}
		if k != strings.ToUpper(k) {
			result.Penalties = append(result.Penalties,
				fmt.Sprintf("%s: lowercase key (-%d)", k, opts.PenaltyLowercaseKey))
			result.Score -= opts.PenaltyLowercaseKey
		}
		if sensitive[k] {
			result.Penalties = append(result.Penalties,
				fmt.Sprintf("%s: sensitive key present (-%d)", k, opts.PenaltySensitiveKey))
			result.Score -= opts.PenaltySensitiveKey
		}
	}

	if result.Score < 0 {
		result.Score = 0
	}
	return result
}
