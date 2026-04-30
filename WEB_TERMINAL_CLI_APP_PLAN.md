# Web Terminal CLI App Plan

## Goal

Build a website-based terminal that lets a person interact with their actual computer from a browser. If they type `pwd`, `ls`, `git status`, or another shell command, it should run on their machine the same way it would in their normal terminal.

The app should also include a friendly side assistant where the user can describe what they want in English, then the app suggests the command or workflow they probably want.

The project should be built with:

- Go for the local backend, terminal process management, filesystem access, and API server.
- JavaScript for the website UI, terminal emulator, and assistant panel.

## Important Design Principle

A normal website cannot directly access a user's hard drive or run real terminal commands. That would be unsafe and browsers block it by design.

To make this work, the user must run a local Go application on their own computer. That Go app exposes a local web server, usually at:

```text
http://localhost:8080
```

The browser UI connects to that local Go server. The Go server runs commands on the user's computer and streams the output back to the website.

## Product Experience

The app should feel like a friendly terminal built into a web application.

Main screen:

- Left or center: real terminal interface.
- Right side: English assistant panel.
- Top bar: current working directory, connection status, and shell selector.
- Bottom/status area: command state, exit code, runtime, and warnings.

Example user flow:

1. User opens the local website.
2. Website connects to the Go backend through WebSocket.
3. User types `pwd`.
4. Go backend runs the command in the current working directory.
5. Output appears in the web terminal.
6. User types in the assistant: "Show me all files, including hidden files, sorted by newest first."
7. Assistant suggests `ls -laht`.
8. User can review, edit, and run the command.

## Architecture

```text
Browser UI
  |
  | HTTP + WebSocket
  v
Local Go Server
  |
  | os/exec + pty
  v
User's Shell
  |
  v
Local Filesystem and Programs
```

## Recommended Project Structure

```text
web-terminal/
  cmd/
    web-terminal/
      main.go
  internal/
    server/
      server.go
      routes.go
      websocket.go
    terminal/
      session.go
      pty.go
      messages.go
    assistant/
      suggest.go
      rules.go
      prompt.go
    safety/
      risk.go
      policy.go
    audit/
      log.go
  web/
    package.json
    index.html
    src/
      main.js
      terminal.js
      assistant.js
      api.js
      state.js
      styles.css
  docs/
    agents/
      README.md
      agent-01-skeleton.md
      agent-02-go-terminal-backend.md
      agent-03-js-terminal-ui.md
      agent-04-assistant-panel.md
      agent-05-safety-audit.md
      agent-06-packaging-qa.md
  README.md
```

## Shared API Contract

All agents should treat this as the initial contract. Agent 1 can refine it during skeleton setup, but later agents should coordinate before changing message shapes.

### WebSocket: Terminal Session

Endpoint:

```text
GET /ws/terminal
```

Client to server messages:

```json
{
  "type": "input",
  "data": "pwd\n"
}
```

```json
{
  "type": "resize",
  "cols": 120,
  "rows": 32
}
```

Server to client messages:

```json
{
  "type": "output",
  "data": "/Users/name/project\n"
}
```

```json
{
  "type": "status",
  "cwd": "/Users/name/project",
  "running": false,
  "shell": "/bin/zsh"
}
```

### HTTP: Assistant Suggestion

Endpoint:

```text
POST /api/assistant/suggest
```

Request:

```json
{
  "prompt": "show me all hidden files sorted by newest first",
  "cwd": "/Users/name/project"
}
```

Response:

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

### HTTP: Command Risk Check

Endpoint:

```text
POST /api/commands/risk
```

Request:

```json
{
  "command": "rm -rf dist"
}
```

Response:

```json
{
  "risk": "high",
  "requiresConfirmation": true,
  "reason": "This command recursively deletes files."
}
```

## Work Split

The project should be split into agent-sized work packages. Agent 1 creates the skeleton first. After that, Agents 2 through 5 can work in parallel against the shared contracts. Agent 6 finishes packaging, integration, and QA.

| Agent | Workstream | Main Ownership | Can Start |
| --- | --- | --- | --- |
| 1 | Skeleton and contracts | Repo layout, Go module, frontend app shell, shared API types | Immediately |
| 2 | Go terminal backend | `internal/server`, `internal/terminal`, PTY WebSocket | After Agent 1 |
| 3 | JavaScript terminal UI | `web/src/terminal.js`, terminal layout, xterm.js integration | After Agent 1 |
| 4 | Assistant panel | `internal/assistant`, `web/src/assistant.js`, suggestion flow | After Agent 1 |
| 5 | Safety and audit | `internal/safety`, `internal/audit`, confirmation policy | After Agent 1 |
| 6 | Packaging and QA | Build scripts, README, integration tests, release flow | After Agents 2-5 |

## Agent Task Files

Use these files as the coding briefs for each person or coding agent:

- [Agent 1: Skeleton and Contracts](docs/agents/agent-01-skeleton.md)
- [Agent 2: Go Terminal Backend](docs/agents/agent-02-go-terminal-backend.md)
- [Agent 3: JavaScript Terminal UI](docs/agents/agent-03-js-terminal-ui.md)
- [Agent 4: Assistant Panel](docs/agents/agent-04-assistant-panel.md)
- [Agent 5: Safety and Audit](docs/agents/agent-05-safety-audit.md)
- [Agent 6: Packaging and QA](docs/agents/agent-06-packaging-qa.md)

## Dependency Flow

```text
Agent 1: Skeleton and contracts
  |
  +--> Agent 2: Go terminal backend
  |
  +--> Agent 3: JavaScript terminal UI
  |
  +--> Agent 4: Assistant panel
  |
  +--> Agent 5: Safety and audit
          |
          v
Agent 6: Packaging, integration, and QA
```

## Core Implementation Choices

### Terminal Process Model

Use a persistent PTY session for the main terminal.

Why:

- Feels closest to a real terminal.
- Supports interactive commands.
- `cd` naturally persists.
- Works well with terminal UIs.

Assistant-generated commands should be previewed and approved before being injected into the PTY.

### Frontend

Use JavaScript for the browser application.

Recommended libraries:

- `xterm.js` for the terminal emulator.
- Vite or a similarly small JS build setup.
- WebSocket client for live terminal input/output.

### Backend

Use Go for the local application server.

Recommended Go packages:

- `net/http` for the local web server.
- `nhooyr.io/websocket` or `gorilla/websocket` for WebSocket communication.
- `github.com/creack/pty` for interactive terminal sessions.

## Safety Requirements

Because this app can run real commands on a real computer, safety is part of the main product.

Required controls:

- Default to manual approval for assistant-generated commands.
- Show command preview before execution.
- Highlight destructive commands.
- Ask for confirmation before deleting files, overwriting files, changing permissions, or running commands with `sudo`.
- Keep a local audit log of commands.
- Bind the server to localhost by default.
- Do not expose the server publicly without explicit advanced configuration.
- Use authentication if remote access is ever added.

## MVP Success Criteria

The MVP is successful when:

- User can start the Go app locally.
- User can open the website at `http://localhost:8080`.
- User can type `pwd` and see the real current directory.
- User can run normal shell commands.
- User can ask the assistant for a command in English.
- Assistant suggests a command without running it automatically.
- User can approve and run the suggested command.
- Risky commands require confirmation.
- Basic command history or audit logging exists.

## Non-Goals For First Version

- Public hosted terminal access.
- Multi-user remote sessions.
- Cloud file syncing.
- Automatic execution of assistant-generated commands.
- Full IDE replacement.

