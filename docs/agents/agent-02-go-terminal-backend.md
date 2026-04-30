# Agent 2: Go Terminal Backend

## Mission

Implement the real local terminal backend. The browser should be able to connect over WebSocket, send terminal input, and receive live output from a shell running on the user's computer.

Start after Agent 1 has created the skeleton.

## Ownership

Primary files and folders:

- `internal/server/websocket.go`
- `internal/server/routes.go`
- `internal/terminal/session.go`
- `internal/terminal/pty.go`
- `internal/terminal/messages.go`

Shared files that may need small edits:

- `cmd/web-terminal/main.go`
- `internal/server/server.go`

## Deliverables

- WebSocket terminal endpoint at `GET /ws/terminal`.
- Persistent PTY-backed shell session.
- Live output streaming from PTY to browser.
- Input forwarding from browser to PTY.
- Terminal resize support.
- Session cleanup when WebSocket disconnects.
- Status messages that include current shell and best-effort current working directory.

## Expected Behavior

User types in the web terminal:

```sh
pwd
```

The backend sends that input into the user's shell and streams the real output back.

## Technical Notes

- Use `github.com/creack/pty` for real terminal behavior.
- Detect the user's shell from `SHELL`; fallback to `/bin/sh` on Unix-like systems.
- Start in the process working directory unless a future setting overrides it.
- Use goroutines carefully:
  - One loop reads WebSocket messages and writes to PTY.
  - One loop reads PTY output and writes WebSocket messages.
  - Both loops must shut down cleanly.
- Do not expose this server beyond localhost by default.

## Message Types

Client to server:

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

Server to client:

```json
{
  "type": "output",
  "data": "..."
}
```

```json
{
  "type": "status",
  "cwd": "/Users/name/project",
  "running": true,
  "shell": "/bin/zsh"
}
```

## Done When

- Browser can connect to `/ws/terminal`.
- `pwd` runs against the real machine.
- `cd` persists because the shell session is persistent.
- Resize messages update PTY size.
- Closing the browser tab cleans up the process.
- Basic tests cover message parsing and session setup where practical.

