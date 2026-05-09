package snapshot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat defines the output format for snapshot exports.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatDotEnv ExportFormat = "dotenv"
)

// Export writes the snapshot variables to w in the specified format.
func Export(s *Snapshot, format ExportFormat, w io.Writer) error {
	switch format {
	case FormatJSON:
		return exportJSON(s, w)
	case FormatCSV:
		return exportCSV(s, w)
	case FormatDotEnv:
		return exportDotEnv(s, w)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportJSON(s *Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s.Vars)
}

func exportCSV(s *Snapshot, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "value"}); err != nil {
		return err
	}
	keys := sortedKeys(s.Vars)
	for _, k := range keys {
		if err := cw.Write([]string{k, s.Vars[k]}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportDotEnv(s *Snapshot, w io.Writer) error {
	keys := sortedKeys(s.Vars)
	for _, k := range keys {
		v := strings.ReplaceAll(s.Vars[k], `"`, `\"`)
		if _, err := fmt.Fprintf(w, "%s=\"%s\"\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
