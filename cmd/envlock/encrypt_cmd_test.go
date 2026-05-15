package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

func setupEncryptStore(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVLOCK_DIR", dir)
	return func() { os.Unsetenv("ENVLOCK_DIR") }
}

func saveEncryptSnapshot(t *testing.T, label string) {
	t.Helper()
	snap := snapshot.Snapshot{
		Label:     label,
		Timestamp: time.Now(),
		Vars:      map[string]string{"APP_ENV": "production", "PORT": "8080"},
	}
	if err := snapshot.DefaultStore().Save(snap); err != nil {
		t.Fatalf("saveEncryptSnapshot: %v", err)
	}
}

func TestRunEncrypt_MissingLabel(t *testing.T) {
	cleanup := setupEncryptStore(t)
	defer cleanup()

	err := runEncrypt([]string{"--passphrase", "secret"}, "encrypt")
	if err == nil {
		t.Fatal("expected error for missing --label")
	}
}

func TestRunEncrypt_MissingPassphrase(t *testing.T) {
	cleanup := setupEncryptStore(t)
	defer cleanup()

	err := runEncrypt([]string{"--label", "snap1"}, "encrypt")
	if err == nil {
		t.Fatal("expected error for missing --passphrase")
	}
}

func TestRunEncrypt_MissingSnapshot(t *testing.T) {
	cleanup := setupEncryptStore(t)
	defer cleanup()

	err := runEncrypt([]string{"--label", "nonexistent", "--passphrase", "secret"}, "encrypt")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunEncrypt_WritesToFile(t *testing.T) {
	cleanup := setupEncryptStore(t)
	defer cleanup()
	saveEncryptSnapshot(t, "snap-enc")

	out := filepath.Join(t.TempDir(), "output.enc")
	err := runEncrypt([]string{"--label", "snap-enc", "--passphrase", "mypassword", "--out", out}, "encrypt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestRunEncrypt_UnknownFlag(t *testing.T) {
	cleanup := setupEncryptStore(t)
	defer cleanup()

	err := runEncrypt([]string{"--bogus", "value"}, "encrypt")
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
