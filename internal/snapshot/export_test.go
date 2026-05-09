package snapshot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func exportSnapshot() *Snapshot {
	return &Snapshot{
		Label:     "test",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Vars: map[string]string{
			"HOME":  "/home/user",
			"GOPATH": "/home/user/go",
			"PATH":  "/usr/bin:/bin",
		},
	}
}

func TestExport_JSON(t *testing.T) {
	s := exportSnapshot()
	var buf bytes.Buffer
	if err := Export(s, FormatJSON, &buf); err != nil {
		t.Fatalf("Export JSON error: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if got["HOME"] != "/home/user" {
		t.Errorf("expected HOME=/home/user, got %q", got["HOME"])
	}
}

func TestExport_CSV(t *testing.T) {
	s := exportSnapshot()
	var buf bytes.Buffer
	if err := Export(s, FormatCSV, &buf); err != nil {
		t.Fatalf("Export CSV error: %v", err)
	}
	output := buf.String()
	if !strings.HasPrefix(output, "key,value") {
		t.Errorf("expected CSV header, got: %q", output)
	}
	if !strings.Contains(output, "HOME,/home/user") {
		t.Errorf("expected HOME row in CSV, got: %q", output)
	}
}

func TestExport_DotEnv(t *testing.T) {
	s := exportSnapshot()
	var buf bytes.Buffer
	if err := Export(s, FormatDotEnv, &buf); err != nil {
		t.Fatalf("Export DotEnv error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, `HOME="/home/user"`) {
		t.Errorf("expected HOME line in dotenv, got: %q", output)
	}
}

func TestExport_DotEnv_EscapesQuotes(t *testing.T) {
	s := &Snapshot{
		Label: "test",
		Vars:  map[string]string{"MSG": `say "hello"`},
	}
	var buf bytes.Buffer
	if err := Export(s, FormatDotEnv, &buf); err != nil {
		t.Fatalf("Export DotEnv error: %v", err)
	}
	if !strings.Contains(buf.String(), `MSG="say \"hello\""`) {
		t.Errorf("quotes not escaped properly: %q", buf.String())
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	s := exportSnapshot()
	var buf bytes.Buffer
	err := Export(s, ExportFormat("xml"), &buf)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
