package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/envlock/internal/snapshot"
)

func setupAnnotateStore(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envlock-annotate-cmd-*")
	if err != nil {
		t.Fatalf("setupAnnotateStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRunAnnotate_MissingSubcommand(t *testing.T) {
	dir := setupAnnotateStore(t)
	var buf bytes.Buffer
	err := runAnnotate([]string{"-store", dir}, &buf)
	if err == nil || !strings.Contains(err.Error(), "subcommand required") {
		t.Errorf("expected subcommand required error, got %v", err)
	}
}

func TestRunAnnotate_SetAndGet(t *testing.T) {
	dir := setupAnnotateStore(t)
	var buf bytes.Buffer
	if err := runAnnotate([]string{"-store", dir, "set", "prod", "production env"}, &buf); err != nil {
		t.Fatalf("set: %v", err)
	}
	if !strings.Contains(buf.String(), "prod") {
		t.Errorf("expected confirmation mentioning label, got %q", buf.String())
	}
	buf.Reset()
	if err := runAnnotate([]string{"-store", dir, "get", "prod"}, &buf); err != nil {
		t.Fatalf("get: %v", err)
	}
	if !strings.Contains(buf.String(), "production env") {
		t.Errorf("expected note in output, got %q", buf.String())
	}
}

func TestRunAnnotate_GetMissing(t *testing.T) {
	dir := setupAnnotateStore(t)
	var buf bytes.Buffer
	err := runAnnotate([]string{"-store", dir, "get", "ghost"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing annotation")
	}
}

func TestRunAnnotate_Remove(t *testing.T) {
	dir := setupAnnotateStore(t)
	store, _ := snapshot.NewAnnotationStore(dir)
	_ = store.Set("staging", "staging note")
	var buf bytes.Buffer
	if err := runAnnotate([]string{"-store", dir, "remove", "staging"}, &buf); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if !strings.Contains(buf.String(), "removed") {
		t.Errorf("expected removed confirmation, got %q", buf.String())
	}
}

func TestRunAnnotate_List(t *testing.T) {
	dir := setupAnnotateStore(t)
	store, _ := snapshot.NewAnnotationStore(dir)
	_ = store.Set("alpha", "note a")
	_ = store.Set("beta", "note b")
	var buf bytes.Buffer
	if err := runAnnotate([]string{"-store", dir, "list"}, &buf); err != nil {
		t.Fatalf("list: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "alpha") || !strings.Contains(output, "beta") {
		t.Errorf("expected both labels in list output, got %q", output)
	}
}

func TestRunAnnotate_ListEmpty(t *testing.T) {
	dir := setupAnnotateStore(t)
	var buf bytes.Buffer
	if err := runAnnotate([]string{"-store", dir, "list"}, &buf); err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if !strings.Contains(buf.String(), "no annotations") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestRunAnnotate_UnknownSubcommand(t *testing.T) {
	dir := setupAnnotateStore(t)
	var buf bytes.Buffer
	err := runAnnotate([]string{"-store", dir, "frobnicate"}, &buf)
	if err == nil || !strings.Contains(err.Error(), "unknown subcommand") {
		t.Errorf("expected unknown subcommand error, got %v", err)
	}
}
