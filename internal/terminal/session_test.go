package terminal

import (
	"strings"
	"testing"
	"time"
)

func TestNewSessionStartsShell(t *testing.T) {
	session, err := NewSession(SessionOptions{
		Shell: fallbackShell,
		Cols:  80,
		Rows:  24,
	})
	if err != nil {
		t.Fatalf("NewSession returned error: %v", err)
	}
	defer session.Close()

	status := session.Status()
	if status.Type != MessageTypeStatus {
		t.Fatalf("status type = %q, want %q", status.Type, MessageTypeStatus)
	}
	if !status.Running {
		t.Fatal("session did not report running shell")
	}
	if status.Shell != fallbackShell {
		t.Fatalf("shell = %q, want %q", status.Shell, fallbackShell)
	}
	if status.CWD == "" {
		t.Fatal("status cwd is empty")
	}
}

func TestSessionRunsInput(t *testing.T) {
	session, err := NewSession(SessionOptions{
		Shell: fallbackShell,
		Cols:  80,
		Rows:  24,
	})
	if err != nil {
		t.Fatalf("NewSession returned error: %v", err)
	}
	defer session.Close()

	if err := session.Write("printf 'agent-two-ready\\n'\n"); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	outputs := make(chan string, 1)
	errors := make(chan error, 1)

	go func() {
		buf := make([]byte, 1024)
		var output strings.Builder
		for {
			n, err := session.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
				if strings.Contains(output.String(), "agent-two-ready") {
					outputs <- output.String()
					return
				}
			}
			if err != nil {
				errors <- err
				return
			}
		}
	}()

	select {
	case <-outputs:
	case err := <-errors:
		t.Fatalf("Read returned error before command output: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("did not receive command output before timeout")
	}
}
