package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultOpenAIBaseURL = "https://api.openai.com/v1"
	defaultOpenAIModel   = "gpt-5.4-mini"
)

var errOpenAINotConfigured = errors.New("openai api key is not configured")

type openAIRequest struct {
	Model           string               `json:"model"`
	Input           []openAIInputMessage `json:"input"`
	MaxOutputTokens int                  `json:"max_output_tokens,omitempty"`
}

type openAIInputMessage struct {
	Role    string                 `json:"role"`
	Content []openAIInputTextBlock `json:"content"`
}

type openAIInputTextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func suggestWithOpenAI(ctx context.Context, req SuggestRequest) (SuggestResponse, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return SuggestResponse{}, errOpenAINotConfigured
	}

	body, err := json.Marshal(openAIRequest{
		Model:           openAIModel(),
		MaxOutputTokens: 700,
		Input: []openAIInputMessage{
			{
				Role: "developer",
				Content: []openAIInputTextBlock{
					{Type: "input_text", Text: systemPrompt},
				},
			},
			{
				Role: "user",
				Content: []openAIInputTextBlock{
					{Type: "input_text", Text: buildOpenAIUserPrompt(req.Input())},
				},
			},
		},
	})
	if err != nil {
		return SuggestResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesURL(), bytes.NewReader(body))
	if err != nil {
		return SuggestResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return SuggestResponse{}, err
	}
	defer resp.Body.Close()

	var decoded openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return SuggestResponse{}, err
	}
	if decoded.Error != nil {
		return SuggestResponse{}, fmt.Errorf("openai response error: %s", decoded.Error.Message)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return SuggestResponse{}, fmt.Errorf("openai response status: %s", resp.Status)
	}

	return parseModelSuggestion(extractOpenAIText(decoded))
}

func openAIModel() string {
	if model := strings.TrimSpace(os.Getenv("OPENAI_MODEL")); model != "" {
		return model
	}
	return defaultOpenAIModel
}

func openAIResponsesURL() string {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")), "/")
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}
	return baseURL + "/responses"
}

func extractOpenAIText(response openAIResponse) string {
	if strings.TrimSpace(response.OutputText) != "" {
		return response.OutputText
	}

	var builder strings.Builder
	for _, item := range response.Output {
		if item.Type != "" && item.Type != "message" {
			continue
		}
		for _, content := range item.Content {
			if content.Type == "" || content.Type == "output_text" || content.Type == "text" {
				builder.WriteString(content.Text)
			}
		}
	}
	return builder.String()
}

func parseModelSuggestion(raw string) (SuggestResponse, error) {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var response SuggestResponse
	if err := json.Unmarshal([]byte(cleaned), &response); err != nil {
		return SuggestResponse{}, err
	}

	response.Suggestions = sanitizeModelSuggestions(response.Suggestions)
	response.Clarification = strings.TrimSpace(response.Clarification)
	if len(response.Suggestions) == 0 && response.Clarification == "" {
		return SuggestResponse{}, errors.New("openai response did not include a suggestion or clarification")
	}
	return response, nil
}

func sanitizeModelSuggestions(suggestions []Suggestion) []Suggestion {
	cleaned := make([]Suggestion, 0, min(len(suggestions), 3))
	for _, suggestion := range suggestions {
		command := strings.TrimSpace(suggestion.Command)
		if command == "" || strings.ContainsAny(command, "\x00\n\r") {
			continue
		}

		risk := suggestion.Risk
		if risk != RiskLow && risk != RiskMedium && risk != RiskHigh {
			risk = RiskMedium
		}

		cleaned = append(cleaned, Suggestion{
			Command:     command,
			Explanation: strings.TrimSpace(suggestion.Explanation),
			Risk:        risk,
		})
		if len(cleaned) == 3 {
			break
		}
	}
	return cleaned
}
