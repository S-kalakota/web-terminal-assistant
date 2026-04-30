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

## Assistant Rules

- Never auto-run suggestions.
- Always show the exact command.
- Keep explanations short.
- Include a risk level: `low`, `medium`, or `high`.
- If the request is ambiguous, return a clarification-style response instead of guessing dangerously.
- For destructive commands, suggest safer alternatives when possible.

## API Response Shape

```json
{
  "suggestions": [
    {
      "command": "ls -laht",
      "explanation": "Lists all files, includes hidden files, uses human-readable sizes, and sorts newest first.",
      "risk": "low"
    }
  ]
}
```

## Integration With Terminal UI

Agent 3 owns the terminal component. This agent should expose or call a simple frontend function such as:

```js
runSuggestedCommand("ls -laht");
```

That function should write the command into the terminal session only after user approval.

## Done When

- User can type English into the assistant panel.
- App displays at least one useful command suggestion.
- User can run an approved suggestion in the terminal.
- Risk level is visible.
- Ambiguous or dangerous requests are handled carefully.

