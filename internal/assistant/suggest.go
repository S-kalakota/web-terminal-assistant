package assistant

import (
	"errors"
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
}

var ErrEmptyInput = errors.New("request body must include text")

func Suggest(req SuggestRequest) (SuggestResponse, error) {
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
	if response, ok := matchAmbiguousRequest(normalized); ok {
		return response, nil
	}

	return clarification("I do not know a safe command for that yet. Try asking for files, git status, your current folder, or a folder creation."), nil
}

func suggestion(command, explanation string, risk RiskLevel) SuggestResponse {
	return SuggestResponse{
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
	}
}
