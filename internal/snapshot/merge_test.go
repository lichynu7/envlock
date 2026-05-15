package snapshot

import (
	"testing"
)

func makeMergeSnapshot(label string, vars map[string]string) Snapshot {
	return Snapshot{Label: label, Vars: vars}
}

func TestMerge_NoConflicts(t *testing.T) {
	base := makeMergeSnapshot("base", map[string]string{"A": "1", "B": "2"})
	overlay := makeMergeSnapshot("overlay", map[string]string{"C": "3"})

	result := Merge(base, overlay, MergeOptions{Strategy: MergeStrategyBase})

	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", result.Conflicts)
	}
	if result.Snapshot.Vars["A"] != "1" || result.Snapshot.Vars["B"] != "2" || result.Snapshot.Vars["C"] != "3" {
		t.Errorf("unexpected merged vars: %v", result.Snapshot.Vars)
	}
}

func TestMerge_ConflictKeepsBase(t *testing.T) {
	base := makeMergeSnapshot("base", map[string]string{"KEY": "base-val"})
	overlay := makeMergeSnapshot("overlay", map[string]string{"KEY": "overlay-val"})

	result := Merge(base, overlay, MergeOptions{Strategy: MergeStrategyBase})

	if len(result.Conflicts) != 1 || result.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", result.Conflicts)
	}
	if result.Snapshot.Vars["KEY"] != "base-val" {
		t.Errorf("expected base value, got %q", result.Snapshot.Vars["KEY"])
	}
}

func TestMerge_ConflictPrefersOverlay(t *testing.T) {
	base := makeMergeSnapshot("base", map[string]string{"KEY": "base-val"})
	overlay := makeMergeSnapshot("overlay", map[string]string{"KEY": "overlay-val"})

	result := Merge(base, overlay, MergeOptions{Strategy: MergeStrategyOverlay})

	if result.Snapshot.Vars["KEY"] != "overlay-val" {
		t.Errorf("expected overlay value, got %q", result.Snapshot.Vars["KEY"])
	}
}

func TestMerge_TracksSides(t *testing.T) {
	base := makeMergeSnapshot("base", map[string]string{"ONLY_BASE": "x"})
	overlay := makeMergeSnapshot("overlay", map[string]string{"ONLY_OVERLAY": "y"})

	result := Merge(base, overlay, MergeOptions{})

	if len(result.BaseOnly) != 1 || result.BaseOnly[0] != "ONLY_BASE" {
		t.Errorf("expected BaseOnly=[ONLY_BASE], got %v", result.BaseOnly)
	}
	if len(result.OverlayOnly) != 1 || result.OverlayOnly[0] != "ONLY_OVERLAY" {
		t.Errorf("expected OverlayOnly=[ONLY_OVERLAY], got %v", result.OverlayOnly)
	}
}

func TestMerge_DefaultLabel(t *testing.T) {
	base := makeMergeSnapshot("snap-a", map[string]string{})
	overlay := makeMergeSnapshot("snap-b", map[string]string{})

	result := Merge(base, overlay, MergeOptions{})

	if result.Snapshot.Label != "merge(snap-a,snap-b)" {
		t.Errorf("unexpected label: %q", result.Snapshot.Label)
	}
}

func TestMerge_CustomLabel(t *testing.T) {
	base := makeMergeSnapshot("a", map[string]string{})
	overlay := makeMergeSnapshot("b", map[string]string{})

	result := Merge(base, overlay, MergeOptions{Label: "my-merge"})

	if result.Snapshot.Label != "my-merge" {
		t.Errorf("expected custom label, got %q", result.Snapshot.Label)
	}
}
