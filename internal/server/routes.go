package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"web-terminal/internal/audit"
	"web-terminal/internal/safety"
	"web-terminal/internal/terminal"
)

// Routes builds the HTTP router for the local app.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/ws/terminal", s.handleTerminalWebSocket)
	mux.HandleFunc("/api/assistant/suggest", s.handleAssistantSuggest)
	mux.HandleFunc("/api/commands/risk", s.handleCommandRisk)
	mux.HandleFunc("/api/commands/audit", s.handleCommandAudit)
	mux.HandleFunc("/", s.handleStatic)

	return securityHeaders(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleAssistantSuggest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var req terminal.AssistantSuggestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "invalid JSON request body",
		})
		return
	}

	prompt := strings.TrimSpace(strings.ToLower(req.Prompt))
	suggestion := terminal.CommandSuggestion{
		Command:     "pwd",
		Explanation: "Placeholder suggestion. The assistant workstream will replace this with real English-to-command behavior.",
		Risk:        terminal.RiskLow,
	}

	if strings.Contains(prompt, "hidden") || strings.Contains(prompt, "newest") {
		suggestion = terminal.CommandSuggestion{
			Command:     "ls -laht",
			Explanation: "Lists all files, includes hidden files, uses human-readable sizes, and sorts newest first.",
			Risk:        terminal.RiskLow,
		}
	}
	suggestion.Risk = safety.AssessCommand(suggestion.Command).Risk

	writeJSON(w, http.StatusOK, terminal.AssistantSuggestResponse{
		Suggestions: []terminal.CommandSuggestion{suggestion},
	})
}

func (s *Server) handleCommandRisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var req terminal.CommandRiskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "invalid JSON request body",
		})
		return
	}

	assessment := safety.AssessCommand(req.Command)
	response := terminal.CommandRiskResponse{
		Risk:                 assessment.Risk,
		RequiresConfirmation: safety.RequiresConfirmation(assessment.Risk),
		Reason:               assessment.Reason,
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleCommandAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var req audit.CommandRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "invalid JSON request body",
		})
		return
	}
	if strings.TrimSpace(req.Command) == "" {
		writeJSON(w, http.StatusBadRequest, terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "command is required",
		})
		return
	}

	assessment := safety.AssessCommand(req.Command)
	req.Risk = assessment.Risk
	req.Source = "assistant"

	if err := audit.AppendCommandRecord(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "failed to write audit log",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		writeMethodNotAllowed(w, http.MethodGet, http.MethodHead)
		return
	}

	path := filepath.Clean(r.URL.Path)
	if path == "/" || path == "." {
		http.ServeFile(w, r, filepath.Join(s.cfg.WebDir, "index.html"))
		return
	}

	relPath := strings.TrimPrefix(path, "/")
	if relPath == ".." || strings.HasPrefix(relPath, "../") {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(s.cfg.WebDir, relPath)
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeMethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeJSON(w, http.StatusMethodNotAllowed, terminal.ErrorMessage{
		Type:    terminal.MessageTypeError,
		Message: "method not allowed",
	})
}
