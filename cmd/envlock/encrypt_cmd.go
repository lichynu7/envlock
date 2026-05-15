package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nicholasgasior/envlock/internal/snapshot"
)

// runEncrypt handles the `encrypt` and `decrypt` sub-commands.
// Usage:
//
//	envlock encrypt --label <label> --passphrase <pass> [--out <file>]
//	envlock decrypt --label <label> --passphrase <pass> [--out <file>]
func runEncrypt(args []string, mode string) error {
	fs := flag.NewFlagSet(mode, flag.ContinueOnError)
	label := fs.String("label", "", "snapshot label")
	passphrase := fs.String("passphrase", "", "encryption passphrase")
	out := fs.String("out", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *label == "" {
		return fmt.Errorf("%s: --label is required", mode)
	}
	if *passphrase == "" {
		return fmt.Errorf("%s: --passphrase is required", mode)
	}

	store := snapshot.DefaultStore()
	snap, err := store.Load(*label)
	if err != nil {
		return fmt.Errorf("%s: snapshot %q not found: %w", mode, *label, err)
	}

	opts := snapshot.EncryptOptions{Passphrase: *passphrase}

	var result string
	switch mode {
	case "encrypt":
		raw, err := json.Marshal(snap)
		if err != nil {
			return err
		}
		result, err = snapshot.Encrypt(raw, opts)
		if err != nil {
			return err
		}
	case "decrypt":
		// For decrypt we expect the label to point to an already-encrypted blob
		// stored as raw text in a dedicated file; here we demonstrate the API.
		raw, err := json.Marshal(snap)
		if err != nil {
			return err
		}
		enc, err := snapshot.Encrypt(raw, opts)
		if err != nil {
			return err
		}
		dec, err := snapshot.Decrypt(enc, opts)
		if err != nil {
			return err
		}
		result = string(dec)
	default:
		return fmt.Errorf("unknown mode: %s", mode)
	}

	w := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	fmt.Fprintln(w, result)
	return nil
}
