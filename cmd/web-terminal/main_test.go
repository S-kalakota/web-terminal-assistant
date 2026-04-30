package main

import "testing"

func TestParseEnvLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantKey   string
		wantValue string
		wantOK    bool
	}{
		{name: "plain", line: "OPENAI_API_KEY=test-key", wantKey: "OPENAI_API_KEY", wantValue: "test-key", wantOK: true},
		{name: "quoted", line: `OPENAI_MODEL="gpt-5.4-mini"`, wantKey: "OPENAI_MODEL", wantValue: "gpt-5.4-mini", wantOK: true},
		{name: "export", line: "export WEB_TERMINAL_ADDR=127.0.0.1:8081", wantKey: "WEB_TERMINAL_ADDR", wantValue: "127.0.0.1:8081", wantOK: true},
		{name: "comment", line: "# ignored", wantOK: false},
		{name: "empty", line: " ", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, value, ok := parseEnvLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if key != tt.wantKey || value != tt.wantValue {
				t.Fatalf("parsed = (%q, %q), want (%q, %q)", key, value, tt.wantKey, tt.wantValue)
			}
		})
	}
}
