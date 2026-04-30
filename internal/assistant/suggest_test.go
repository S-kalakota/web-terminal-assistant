package assistant

import (
	"errors"
	"testing"
)

func TestSuggestMVPExamples(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		command string
	}{
		{name: "where am I", text: "where am I", command: "pwd"},
		{name: "list files", text: "list files", command: "ls"},
		{name: "show hidden files", text: "show hidden files", command: "ls -la"},
		{name: "show files sorted by newest", text: "show files sorted by newest", command: "ls -laht"},
		{name: "what changed in git", text: "what changed in git", command: "git status"},
		{name: "show current branch", text: "show current branch", command: "git branch --show-current"},
		{name: "make folder", text: "make a folder named reports", command: "mkdir reports"},
		{name: "find large files", text: "find large files", command: "find . -type f -size +100M"},
		{name: "print working directory", text: "print working directory", command: "pwd"},
		{name: "show directory contents", text: "show directory contents", command: "ls"},
		{name: "show git status", text: "show git status", command: "git status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := Suggest(SuggestRequest{Text: tt.text})
			if err != nil {
				t.Fatalf("Suggest returned error: %v", err)
			}
			if len(response.Suggestions) != 1 {
				t.Fatalf("Suggestions length = %d, want 1", len(response.Suggestions))
			}
			if response.Suggestions[0].Command != tt.command {
				t.Fatalf("Command = %q, want %q", response.Suggestions[0].Command, tt.command)
			}
			if response.Suggestions[0].Risk != RiskLow {
				t.Fatalf("Risk = %q, want %q", response.Suggestions[0].Risk, RiskLow)
			}
		})
	}
}

func TestSuggestAcceptsPromptForCompatibility(t *testing.T) {
	response, err := Suggest(SuggestRequest{Prompt: "list files"})
	if err != nil {
		t.Fatalf("Suggest returned error: %v", err)
	}
	if response.Suggestions[0].Command != "ls" {
		t.Fatalf("Command = %q, want ls", response.Suggestions[0].Command)
	}
}

func TestSuggestRejectsEmptyInput(t *testing.T) {
	_, err := Suggest(SuggestRequest{Text: "  "})
	if !errors.Is(err, ErrEmptyInput) {
		t.Fatalf("error = %v, want ErrEmptyInput", err)
	}
}

func TestSuggestClarifiesAmbiguousInput(t *testing.T) {
	tests := []string{
		"delete stuff",
		"clean this project",
		"make a folder",
		"fix git",
		"install dependencies",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			response, err := Suggest(SuggestRequest{Text: text})
			if err != nil {
				t.Fatalf("Suggest returned error: %v", err)
			}
			if len(response.Suggestions) != 0 {
				t.Fatalf("Suggestions length = %d, want 0", len(response.Suggestions))
			}
			if response.Clarification == "" {
				t.Fatal("Clarification is empty")
			}
		})
	}
}

func TestSuggestQuotesFolderNames(t *testing.T) {
	response, err := Suggest(SuggestRequest{Text: "make a folder named Project Notes"})
	if err != nil {
		t.Fatalf("Suggest returned error: %v", err)
	}
	if response.Suggestions[0].Command != "mkdir 'Project Notes'" {
		t.Fatalf("Command = %q, want quoted folder name", response.Suggestions[0].Command)
	}
}

func TestSuggestClarifiesUnsafeFolderNames(t *testing.T) {
	response, err := Suggest(SuggestRequest{Text: "make a folder named foo; rm -rf ."})
	if err != nil {
		t.Fatalf("Suggest returned error: %v", err)
	}
	if len(response.Suggestions) != 0 {
		t.Fatalf("Suggestions length = %d, want 0", len(response.Suggestions))
	}
	if response.Clarification == "" {
		t.Fatal("Clarification is empty")
	}
}
