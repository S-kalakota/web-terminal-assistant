package server

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"web-terminal/internal/terminal"
)

func TestTerminalWebSocketRunsShellInput(t *testing.T) {
	t.Setenv("SHELL", "/bin/sh")

	app := New(Config{})
	server := httptest.NewServer(app.Routes())
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/terminal"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial returned error: %v", err)
	}
	defer conn.Close()

	var status terminal.TerminalStatusMessage
	if err := conn.ReadJSON(&status); err != nil {
		t.Fatalf("ReadJSON status returned error: %v", err)
	}
	if status.Type != terminal.MessageTypeStatus {
		t.Fatalf("initial message type = %q, want %q", status.Type, terminal.MessageTypeStatus)
	}
	if status.Shell != "/bin/sh" {
		t.Fatalf("shell = %q, want /bin/sh", status.Shell)
	}

	if err := conn.WriteJSON(terminal.TerminalInputMessage{
		Type: terminal.MessageTypeInput,
		Data: "printf 'websocket-ready\\n'\n",
	}); err != nil {
		t.Fatalf("WriteJSON input returned error: %v", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline returned error: %v", err)
	}

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("ReadMessage returned error before command output: %v", err)
		}

		var output terminal.TerminalOutputMessage
		if err := json.Unmarshal(payload, &output); err != nil {
			t.Fatalf("Unmarshal output returned error: %v", err)
		}
		if output.Type == terminal.MessageTypeOutput && strings.Contains(output.Data, "websocket-ready") {
			return
		}
	}
}
