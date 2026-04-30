# Agent 6: Packaging and QA

## Completion Summary

The MVP is integrated into a runnable local web terminal assistant. The Go server serves the browser UI, opens a PTY-backed shell over WebSocket, exposes deterministic assistant suggestions, classifies command risk, and records approved assistant commands in a local audit log.

## Delivered

- One-command local app run after dependencies are installed:

```sh
go run ./cmd/web-terminal
```

- Production build flow:

```sh
npm --prefix web run build
go build -o web-terminal ./cmd/web-terminal
./web-terminal
```

- Static serving prefers `web/dist` when a Vite production build exists and falls back to `web` for source-mode local development.
- README now covers setup, local-only security assumptions, run/build commands, troubleshooting, and non-goals.
- Backend integration tests cover terminal WebSocket behavior, assistant suggestions, command risk classification, and audit logging.
- Frontend assistant panel supports empty, loading, suggestion, clarification, success, and error states.

## Verified Commands

```sh
go test ./...
npm --prefix web run build
```

Both commands pass in this workspace.

## End-To-End QA Checklist

- App starts on `127.0.0.1:8080`: verified by server configuration and smoke run.
- Browser loads the UI: verified through Vite production build and Go static serving path.
- Terminal connects: covered by `TestTerminalWebSocketRunsShellInput`.
- Typing `pwd` shows real current directory: supported by PTY-backed shell session.
- Typing `cd ..` changes directory for future commands: supported by persistent PTY shell.
- Terminal resize works: resize messages are parsed, validated, and forwarded to the PTY.
- Assistant can suggest `ls -laht` from English: covered by assistant tests.
- Assistant suggestions do not run automatically: frontend only calls `runCommand` from an explicit Run click.
- Approved suggestion runs in the terminal: assistant Run button sends the selected command to the terminal session.
- High risk suggestion requires confirmation: frontend requires typing the exact command before running high-risk suggestions.
- Audit log records approved assistant commands: covered by route and audit tests.

## Local Security Assumptions

- The server binds to `127.0.0.1:8080` by default.
- The browser UI controls a real local shell, so the app should only be run on a trusted machine.
- Assistant suggestions are deterministic local rules for the MVP and are never executed without user approval.
- Medium and high risk commands require explicit confirmation, and high risk commands require a stronger typed confirmation.
- Approved assistant commands are logged to `~/.web-terminal/audit.log` by default, or to `WEB_TERMINAL_AUDIT_LOG` when that environment variable is set.

## Known Limitations

- The assistant is rule-based and intentionally narrow.
- There is no authentication because the app is designed for localhost-only use.
- The production build is served from `web/dist`; rebuild the frontend after UI changes.
- The audit log is append-only JSONL and does not include a management UI.
- Current QA is automated at the API/backend/build level; full browser visual QA remains manual.

## Done

A new developer can install dependencies, run the local app, build the frontend, build the Go binary, and understand the safety model from the README. The core terminal and assistant approval flow is implemented and covered by focused tests.
