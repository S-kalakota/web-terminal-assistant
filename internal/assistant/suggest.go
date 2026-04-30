package assistant

import (
	"context"
	"errors"
	"log"
	"strings"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type SuggestRequest struct {
	Text   string `json:"text"`
	Prompt string `json:"prompt,omitempty"`
	CWD    string `json:"cwd,omitempty"`
}

func (r SuggestRequest) Input() string {
	if strings.TrimSpace(r.Text) != "" {
		return r.Text
	}
	return r.Prompt
}

type Suggestion struct {
	Command     string    `json:"command"`
	Explanation string    `json:"explanation"`
	Risk        RiskLevel `json:"risk"`
}

type SuggestResponse struct {
	Suggestions   []Suggestion `json:"suggestions"`
	Clarification string       `json:"clarification,omitempty"`
	Error         string       `json:"error,omitempty"`
	Source        string       `json:"source,omitempty"`
	Warning       string       `json:"warning,omitempty"`
}

var ErrEmptyInput = errors.New("request body must include text")

func Suggest(req SuggestRequest) (SuggestResponse, error) {
	return suggestRules(req)
}

func SuggestWithLLM(ctx context.Context, req SuggestRequest) (SuggestResponse, error) {
	input := strings.TrimSpace(req.Input())
	if input == "" {
		return SuggestResponse{}, ErrEmptyInput
	}

	if response, err := suggestWithOpenAI(ctx, req); err == nil {
		response.Source = "openai"
		return response, nil
	} else if !errors.Is(err, errOpenAINotConfigured) {
		log.Printf("assistant OpenAI unavailable, using local rules: %v", err)
		response, fallbackErr := suggestRules(req)
		response.Warning = "OpenAI was unavailable, so the assistant used local fallback rules."
		return response, fallbackErr
	}

	return suggestRules(req)
}

func suggestRules(req SuggestRequest) (SuggestResponse, error) {
	input := strings.TrimSpace(req.Input())
	if input == "" {
		return SuggestResponse{}, ErrEmptyInput
	}

	normalized := normalize(input)
	if response, ok := matchKnownRequest(normalized); ok {
		return response, nil
	}
	if response, ok := matchFolderCreation(input, normalized); ok {
		return response, nil
	}
	if response, ok := matchDirectoryNavigation(input, normalized); ok {
		return response, nil
	}
	if response, ok := matchAmbiguousRequest(normalized); ok {
		return response, nil
	}
	if response, ok := fallbackSuggestion(normalized); ok {
		return response, nil
	}

	return suggestion("pwd", "Prints the current directory so you can orient before the next command.", RiskLow), nil
}

func suggestion(command, explanation string, risk RiskLevel) SuggestResponse {
	return SuggestResponse{
		Source: "rules",
		Suggestions: []Suggestion{
			{
				Command:     command,
				Explanation: explanation,
				Risk:        risk,
			},
		},
	}
}

func clarification(message string) SuggestResponse {
	return SuggestResponse{
		Suggestions:   []Suggestion{},
		Clarification: message,
		Source:        "rules",
	}
}
