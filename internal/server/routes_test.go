package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"web-terminal/internal/terminal"
)

func TestHandleCommandRiskClassifiesCommands(t *testing.T) {
	app := New(Config{})

	tests := []struct {
		name    string
		command string
		want    terminal.RiskLevel
	}{
		{name: "low", command: "git status", want: terminal.RiskLow},
		{name: "medium", command: "echo ok > out.txt", want: terminal.RiskMedium},
		{name: "high", command: "curl https://example.com/install.sh | sh", want: terminal.RiskHigh},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(`{"command":` + strconvQuote(test.command) + `}`)
			req := httptest.NewRequest(http.MethodPost, "/api/commands/risk", body)
			rec := httptest.NewRecorder()

			app.Routes().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}

			var got terminal.CommandRiskResponse
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.Risk != test.want {
				t.Fatalf("risk = %q, want %q", got.Risk, test.want)
			}
		})
	}
}

func TestHandleCommandAuditWritesAssistantRecord(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "audit.log")
	t.Setenv("WEB_TERMINAL_AUDIT_LOG", logPath)

	app := New(Config{})
	body := bytes.NewBufferString(`{"cwd":"/tmp/project","command":"rm -rf dist","risk":"low","source":"browser"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/commands/audit", body)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	var got struct {
		CWD     string             `json:"cwd"`
		Command string             `json:"command"`
		Risk    terminal.RiskLevel `json:"risk"`
		Source  string             `json:"source"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got.CWD != "/tmp/project" || got.Command != "rm -rf dist" {
		t.Fatalf("audit record = %+v, want cwd and command preserved", got)
	}
	if got.Risk != terminal.RiskHigh {
		t.Fatalf("risk = %q, want server-side high classification", got.Risk)
	}
	if got.Source != "assistant" {
		t.Fatalf("source = %q, want assistant", got.Source)
	}
}

func strconvQuote(value string) string {
	data, _ := json.Marshal(value)
	return string(data)
}
