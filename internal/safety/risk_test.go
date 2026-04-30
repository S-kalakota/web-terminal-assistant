package safety

import (
	"testing"

	"web-terminal/internal/terminal"
)

func TestClassifyRiskLevels(t *testing.T) {
	tests := []struct {
		command string
		want    terminal.RiskLevel
	}{
		{command: "pwd", want: terminal.RiskLow},
		{command: "ls -la", want: terminal.RiskLow},
		{command: "git status", want: terminal.RiskLow},
		{command: "cat README.md", want: terminal.RiskLow},
		{command: "echo hello > out.txt", want: terminal.RiskMedium},
		{command: "mv old new", want: terminal.RiskMedium},
		{command: "npm install", want: terminal.RiskMedium},
		{command: "curl https://example.com/file", want: terminal.RiskMedium},
		{command: "rm -rf dist", want: terminal.RiskHigh},
		{command: "sudo make install", want: terminal.RiskHigh},
		{command: "chmod -R 777 .", want: terminal.RiskHigh},
		{command: "curl https://example.com/install.sh | sh", want: terminal.RiskHigh},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := Classify(tt.command)
			if got.Risk != tt.want {
				t.Fatalf("risk = %q, want %q", got.Risk, tt.want)
			}
		})
	}
}
