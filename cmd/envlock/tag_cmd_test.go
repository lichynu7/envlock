package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func setupTagStore(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envlock-tag-cmd-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRunTag_AddAndList(t *testing.T) {
	dir := setupTagStore(t)
	var out bytes.Buffer

	if err := runTag([]string{"release", "snap-1"}, &out, dir); err != nil {
		t.Fatalf("runTag add: %v", err)
	}
	out.Reset()

	if err := runTag([]string{"--list", "release"}, &out, dir); err != nil {
		t.Fatalf("runTag list: %v", err)
	}
	if !strings.Contains(out.String(), "snap-1") {
		t.Errorf("expected snap-1 in output, got: %s", out.String())
	}
}

func TestRunTag_ListAllTags(t *testing.T) {
	dir := setupTagStore(t)
	var out bytes.Buffer

	_ = runTag([]string{"alpha", "snap-a"}, &out, dir)
	_ = runTag([]string{"beta", "snap-b"}, &out, dir)
	out.Reset()

	if err := runTag([]string{"--list"}, &out, dir); err != nil {
		t.Fatalf("runTag list all: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "alpha") || !strings.Contains(result, "beta") {
		t.Errorf("expected both tags in output, got: %s", result)
	}
}

func TestRunTag_Remove(t *testing.T) {
	dir := setupTagStore(t)
	var out bytes.Buffer

	_ = runTag([]string{"ci", "snap-1"}, &out, dir)
	out.Reset()

	if err := runTag([]string{"--remove", "ci", "snap-1"}, &out, dir); err != nil {
		t.Fatalf("runTag remove: %v", err)
	}
	if !strings.Contains(out.String(), "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out.String())
	}
}

func TestRunTag_MissingArgs(t *testing.T) {
	dir := setupTagStore(t)
	var out bytes.Buffer
	err := runTag([]string{"onlyone"}, &out, dir)
	if err == nil {
		t.Error("expected error for missing label argument")
	}
}

func TestRunTag_ListEmpty(t *testing.T) {
	dir := setupTagStore(t)
	var out bytes.Buffer
	if err := runTag([]string{"--list"}, &out, dir); err != nil {
		t.Fatalf("runTag list empty: %v", err)
	}
	if !strings.Contains(out.String(), "no tags") {
		t.Errorf("expected 'no tags' message, got: %s", out.String())
	}
}
