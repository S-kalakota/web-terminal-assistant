package terminal

import "testing"

func TestParseClientMessageInput(t *testing.T) {
	msg, err := ParseClientMessage([]byte(`{"type":"input","data":"pwd\n"}`))
	if err != nil {
		t.Fatalf("ParseClientMessage returned error: %v", err)
	}

	if msg.Type != MessageTypeInput {
		t.Fatalf("Type = %q, want %q", msg.Type, MessageTypeInput)
	}
	if msg.Data != "pwd\n" {
		t.Fatalf("Data = %q, want pwd newline", msg.Data)
	}
}

func TestParseClientMessageResize(t *testing.T) {
	msg, err := ParseClientMessage([]byte(`{"type":"resize","cols":120,"rows":32}`))
	if err != nil {
		t.Fatalf("ParseClientMessage returned error: %v", err)
	}

	if msg.Cols != 120 || msg.Rows != 32 {
		t.Fatalf("resize = %dx%d, want 120x32", msg.Cols, msg.Rows)
	}
}

func TestParseClientMessageRejectsBadResize(t *testing.T) {
	if _, err := ParseClientMessage([]byte(`{"type":"resize","cols":0,"rows":32}`)); err == nil {
		t.Fatal("ParseClientMessage accepted resize with zero cols")
	}
}

func TestParseClientMessageRejectsUnknownType(t *testing.T) {
	if _, err := ParseClientMessage([]byte(`{"type":"status"}`)); err == nil {
		t.Fatal("ParseClientMessage accepted unsupported client message type")
	}
}
