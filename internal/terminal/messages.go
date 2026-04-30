package terminal

import (
	"encoding/json"
	"errors"
	"fmt"
)

// MessageType identifies messages exchanged between the browser and backend.
type MessageType string

const (
	MessageTypeInput  MessageType = "input"
	MessageTypeResize MessageType = "resize"
	MessageTypeOutput MessageType = "output"
	MessageTypeStatus MessageType = "status"
	MessageTypeError  MessageType = "error"
)

// RiskLevel describes how careful the UI should be before running a command.
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// TerminalInputMessage is sent by the browser when the user types.
type TerminalInputMessage struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

// TerminalResizeMessage is sent by the browser when terminal dimensions change.
type TerminalResizeMessage struct {
	Type MessageType `json:"type"`
	Cols int         `json:"cols"`
	Rows int         `json:"rows"`
}

// TerminalClientMessage is the normalized form of browser-to-backend websocket messages.
type TerminalClientMessage struct {
	Type MessageType `json:"type"`
	Data string      `json:"data,omitempty"`
	Cols int         `json:"cols,omitempty"`
	Rows int         `json:"rows,omitempty"`
}

// TerminalOutputMessage is sent by the backend with shell output.
type TerminalOutputMessage struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

// TerminalStatusMessage is sent by the backend with session state.
type TerminalStatusMessage struct {
	Type    MessageType `json:"type"`
	CWD     string      `json:"cwd"`
	Running bool        `json:"running"`
	Shell   string      `json:"shell"`
}

// ErrorMessage is a structured API or websocket error payload.
type ErrorMessage struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// ParseClientMessage validates a browser websocket message before it touches the PTY.
func ParseClientMessage(payload []byte) (TerminalClientMessage, error) {
	var msg TerminalClientMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return TerminalClientMessage{}, fmt.Errorf("invalid JSON message: %w", err)
	}

	switch msg.Type {
	case MessageTypeInput:
		return msg, nil
	case MessageTypeResize:
		if msg.Cols <= 0 || msg.Rows <= 0 {
			return TerminalClientMessage{}, errors.New("resize messages require positive cols and rows")
		}
		return msg, nil
	default:
		return TerminalClientMessage{}, fmt.Errorf("unsupported terminal message type %q", msg.Type)
	}
}

// AssistantSuggestRequest asks the backend to translate English to commands.
type AssistantSuggestRequest struct {
	Prompt string `json:"prompt"`
	CWD    string `json:"cwd"`
}

// AssistantSuggestResponse returns possible commands for a user request.
type AssistantSuggestResponse struct {
	Suggestions []CommandSuggestion `json:"suggestions"`
}

// CommandSuggestion is a single assistant-generated command candidate.
type CommandSuggestion struct {
	Command     string    `json:"command"`
	Explanation string    `json:"explanation"`
	Risk        RiskLevel `json:"risk"`
}

// CommandRiskRequest asks the backend to classify a command.
type CommandRiskRequest struct {
	Command string `json:"command"`
}

// CommandRiskResponse describes the risk and confirmation requirement.
type CommandRiskResponse struct {
	Risk                 RiskLevel `json:"risk"`
	RequiresConfirmation bool      `json:"requiresConfirmation"`
	Reason               string    `json:"reason"`
}
