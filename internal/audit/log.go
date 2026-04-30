package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"web-terminal/internal/terminal"
)

// Record captures one approved assistant command.
type Record struct {
	Timestamp time.Time          `json:"timestamp"`
	CWD       string             `json:"cwd"`
	Command   string             `json:"command"`
	Risk      terminal.RiskLevel `json:"risk"`
	Source    string             `json:"source"`
}

// CommandRecord captures one approved assistant command from the HTTP API.
type CommandRecord = Record

// Append writes a command approval event as one JSON line in the local audit log.
func Append(record Record) error {
	path, err := DefaultPath()
	if err != nil {
		return err
	}

	return AppendTo(path, record)
}

// AppendCommandRecord writes a command approval event as one JSON line in the local audit log.
func AppendCommandRecord(record CommandRecord) error {
	return Append(record)
}

// AppendTo writes a command approval event to a specific audit log path.
func AppendTo(path string, record Record) error {
	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now().UTC()
	}
	if record.Source == "" {
		record.Source = "assistant"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoded, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return err
	}

	return nil
}

// DefaultPath returns ~/.web-terminal/audit.log.
func DefaultPath() (string, error) {
	if path := os.Getenv("WEB_TERMINAL_AUDIT_LOG"); path != "" {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".web-terminal", "audit.log"), nil
}
