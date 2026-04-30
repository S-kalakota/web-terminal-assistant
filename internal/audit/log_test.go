package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"web-terminal/internal/terminal"
)

func TestAppendToWritesJSONLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")

	err := AppendTo(path, Record{
		Timestamp: time.Date(2026, 4, 29, 21, 30, 0, 0, time.UTC),
		CWD:       "/tmp/project",
		Command:   "ls -laht",
		Risk:      terminal.RiskLow,
		Source:    "assistant",
	})
	if err != nil {
		t.Fatalf("AppendTo returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if record.Command != "ls -laht" {
		t.Fatalf("command = %q, want ls -laht", record.Command)
	}
}
