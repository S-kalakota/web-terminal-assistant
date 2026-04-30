# Agent 4: Assistant Panel

## Mission

Build the English assistant panel that translates user intent into suggested shell commands. The assistant should suggest commands, explain them, and let the user approve running them in the terminal.

Start after Agent 1 has created the skeleton.

## Ownership

Primary files and folders:

- `internal/assistant/suggest.go`
- `internal/assistant/rules.go`
- `internal/assistant/prompt.go`
- `web/src/assistant.js`
- `web/src/api.js`
- `web/src/styles.css`

Shared files that may need small edits:

- `internal/server/routes.go`
- `web/src/main.js`

## Deliverables

- Assistant panel UI with English input.
- `POST /api/assistant/suggest` endpoint.
- Rule-based MVP suggestions for common tasks.
- Suggested command preview.
- Explanation text.
- Risk level display.
- Run button that sends the approved command to the terminal UI.
- Empty, loading, success, clarification, and error states.
- Focus and keyboard behavior that works without a mouse.

## MVP Suggestion Examples

Support at least these English requests:

- "where am I" -> `pwd`
- "list files" -> `ls`
- "show hidden files" -> `ls -la`
- "show files sorted by newest" -> `ls -laht`
- "what changed in git" -> `git status`
- "show current branch" -> `git branch --show-current`
- "make a folder named X" -> `mkdir X`
- "find large files" -> `find . -type f -size +100M`
- "print working directory" -> `pwd`
- "show directory contents" -> `ls`
- "show git status" -> `git status`

The implementation can support more phrases, but the examples above are the minimum acceptance set.

## Assistant Rules

- Never auto-run suggestions.
- Always show the exact command.
- Keep explanations short.
- Include a risk level: `low`, `medium`, or `high`.
- If the request is ambiguous, return a clarification-style response instead of guessing dangerously.
- For destructive commands, suggest safer alternatives when possible.
- Prefer read-only commands for the MVP.
- Do not call external AI services for this phase. The MVP should be deterministic and local.
- Do not generate shell syntax that depends on untrusted command substitution.

## API Contract

```text
POST /api/assistant/suggest
```

Request:

```json
{
  "text": "show files sorted by newest",
  "cwd": "/Users/name/project"
}
```

The `cwd` field is optional and should be treated as context only. The assistant should not fail if it is absent.

Successful response:

```json
{
  "suggestions": [
    {
      "command": "ls -laht",
      "explanation": "Lists all files, includes hidden files, uses human-readable sizes, and sorts newest first.",
      "risk": "low"
    }
  ],
  "clarification": ""
}
```

Clarification response:

```json
{
  "suggestions": [],
  "clarification": "Which folder should I create?"
}
```

Error response:

```json
{
  "error": "Request body must include text."
}
```

## Backend Requirements

- Add an `internal/assistant` package with small, testable functions.
- Keep request matching in `rules.go`.
- Keep public suggestion orchestration in `suggest.go`.
- Keep any prompt or wording templates in `prompt.go`, even if the MVP is not using an AI model yet.
- Normalize input by trimming whitespace, lowercasing for matching, and collapsing repeated spaces.
- Preserve user-provided names when building commands such as `mkdir X`.
- Quote user-provided paths or names safely when needed.
- Return no more than three suggestions for one request.
- Return HTTP `400` for malformed JSON or empty input.

Suggested Go types:

```go
type SuggestRequest struct {
	Text string `json:"text"`
	CWD  string `json:"cwd,omitempty"`
}

type Suggestion struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
	Risk        string `json:"risk"`
}

type SuggestResponse struct {
	Suggestions   []Suggestion `json:"suggestions"`
	Clarification string       `json:"clarification,omitempty"`
	Error         string       `json:"error,omitempty"`
}
```

## Rule Matching Notes

Start simple and deterministic. Good matching is more important than broad matching.

Recommended rule order:

1. Exact or near-exact known read-only requests.
2. Folder creation requests with an extracted name.
3. Known safe search or git requests.
4. Ambiguous requests that need clarification.
5. Unsupported requests with a polite clarification.

Examples that should ask for clarification instead of guessing:

- "delete stuff"
- "clean this project"
- "make a folder"
- "fix git"
- "install dependencies"

Examples that should avoid direct destructive suggestions:

- "delete all node modules" should explain the risk or ask for confirmation context rather than returning `rm -rf node_modules`.
- "reset git" should ask what kind of reset the user means.
- "remove everything" should not produce a runnable command.

## Frontend Requirements

- Add a real assistant panel in `web/src/assistant.js`.
- Use `web/src/api.js` for the HTTP request.
- Render suggestions with command preview, explanation, risk label, and Run button.
- Disable controls while a suggestion request is loading.
- Show errors without clearing the user's typed request.
- Keep the assistant usable on narrow screens.
- Do not send a command to the terminal until the user clicks Run.
- After Run, keep the suggestion visible so the user can see what was approved.

## Integration With Terminal UI

Agent 3 owns the terminal component. This agent should expose or call a simple frontend function such as:

```js
runSuggestedCommand("ls -laht");
```

That function should write the command into the terminal session only after user approval. If Agent 3 has already created a different terminal API, adapt to it rather than replacing terminal behavior.

The assistant should send a trailing newline only if the terminal API expects executable input. If the terminal API supports staging text before execution, prefer staging the command and let the terminal module own execution details.

## Testing Requirements

Backend tests should cover:

- Each MVP suggestion example.
- Empty input.
- Ambiguous input.
- Folder name extraction.
- Basic shell quoting for generated folder names.

Frontend testing can be manual for the MVP, but verify:

- Suggest button calls the API.
- Loading state appears.
- Suggestions render correctly.
- Clarifications render without Run buttons.
- Clicking Run sends only the approved command to the terminal.

## Non-Goals

- Do not build an LLM integration in this agent.
- Do not implement full natural language understanding.
- Do not build the risk audit log. Agent 5 owns audit logging.
- Do not replace the terminal WebSocket behavior. Agent 2 owns terminal backend behavior.

## Done When

- User can type English into the assistant panel.
- App displays at least one useful command suggestion.
- User can run an approved suggestion in the terminal.
- Risk level is visible.
- Ambiguous or dangerous requests are handled carefully.
- `POST /api/assistant/suggest` returns the documented JSON shapes.
- Tests cover the rule-based backend behavior.
- The assistant panel works with Agent 3's terminal UI without breaking direct terminal input.
