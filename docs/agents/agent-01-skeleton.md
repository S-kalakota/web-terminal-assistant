# Agent 1: Skeleton and Contracts

## Mission

Create the initial Go and JavaScript project skeleton so the rest of the team can build against stable files, routes, and message contracts.

This agent goes first.

## Ownership

Primary files and folders:

- `go.mod`
- `cmd/web-terminal/main.go`
- `internal/server/server.go`
- `internal/server/routes.go`
- `internal/terminal/messages.go`
- `web/package.json`
- `web/index.html`
- `web/src/main.js`
- `web/src/api.js`
- `web/src/state.js`
- `web/src/styles.css`
- `README.md`

## Deliverables

- A Go module that builds.
- A local HTTP server that binds to `127.0.0.1:8080` by default.
- Static frontend serving from the Go server.
- Placeholder routes for:
  - `GET /`
  - `GET /healthz`
  - `GET /ws/terminal`
  - `POST /api/assistant/suggest`
  - `POST /api/commands/risk`
- Shared request and response structs for terminal, assistant, and risk messages.
- Minimal JS app shell that loads in the browser.
- README with local development commands.

## API Contract To Establish

### Health Check

```text
GET /healthz
```

Response:

```json
{
  "ok": true
}
```

### Terminal WebSocket

```text
GET /ws/terminal
```

Message types:

- `input`
- `resize`
- `output`
- `status`
- `error`

### Assistant Suggestion

```text
POST /api/assistant/suggest
```

### Command Risk Check

```text
POST /api/commands/risk
```

## Implementation Notes

- Keep placeholder handlers simple and clearly marked.
- Use JSON structs that later agents can import instead of duplicating shape definitions.
- Use normal Go package boundaries under `internal/`.
- Set the default server address to `127.0.0.1:8080`.
- Avoid adding real PTY behavior here. Agent 2 owns that.
- Avoid building the final UI here. Agent 3 owns that.

## Done When

- `go run ./cmd/web-terminal` starts the server.
- Browser can load the app.
- `GET /healthz` returns JSON.
- Placeholder assistant and risk endpoints return valid JSON.
- The project has clear folders for the next agents.

