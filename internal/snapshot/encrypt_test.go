package snapshot

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	opts := EncryptOptions{Passphrase: "supersecret"}
	plaintext := []byte("HOME=/home/user\nPATH=/usr/bin")

	encoded, err := Encrypt(plaintext, opts)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if encoded == "" {
		t.Fatal("expected non-empty encoded string")
	}

	got, err := Decrypt(encoded, opts)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("round-trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := Encrypt([]byte("data"), EncryptOptions{})
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	_, err := Decrypt("somedata", EncryptOptions{})
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encoded, err := Encrypt([]byte("secret"), EncryptOptions{Passphrase: "correct"})
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	_, err = Decrypt(encoded, EncryptOptions{Passphrase: "wrong"})
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestEncrypt_ProducesBase64(t *testing.T) {
	opts := EncryptOptions{Passphrase: "testkey"}
	encoded, err := Encrypt([]byte("VALUE=1"), opts)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	// base64 standard encoding only contains A-Z a-z 0-9 + / =
	for _, ch := range encoded {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", ch) {
			t.Errorf("unexpected character in base64 output: %q", ch)
		}
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("!!!notbase64!!!", EncryptOptions{Passphrase: "key"})
	if err == nil {
		t.Fatal("expected error for invalid base64 input")
	}
}

func TestDefaultEncryptOptions(t *testing.T) {
	opts := DefaultEncryptOptions()
	if opts.Passphrase != "" {
		t.Errorf("expected empty passphrase, got %q", opts.Passphrase)
	}
}
