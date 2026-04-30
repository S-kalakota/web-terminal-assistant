# Web Terminal Assistant

Web Terminal Assistant is a local-first web app that gives the browser a real terminal session on your machine. A Go server runs locally, serves the UI, opens a PTY-backed shell, and streams terminal input/output over WebSocket.

The app also includes a command assistant. You describe a task in English, review the suggested command, and choose whether to run it. Suggestions never execute automatically.

## Why It Runs Locally

Browsers are not allowed to run shell commands or freely access your filesystem. This project works by running a local Go server on your computer and binding to:

```text
127.0.0.1:8080
```

The browser talks only to that local server. Because the server controls a real shell, only run it on a trusted machine.

## Features

- Browser terminal powered by xterm.js.
- PTY-backed local shell with persistent session state.
- WebSocket terminal input, output, and resize handling.
- Rule-based English-to-command suggestions for common safe tasks.
- Optional OpenAI-powered assistant suggestions through the Responses API.
- Risk classification for low, medium, and high risk commands.
- Extra confirmation for high risk assistant commands.
- Local audit log for approved assistant commands.

## Requirements

- Go 1.22 or newer.
- Node.js and npm.
- A Unix-like shell environment.

## Install

From the project root:

```sh
npm --prefix web install
go mod download
```

## Run Locally

```sh
go run ./cmd/web-terminal
```

Open:

```text
http://127.0.0.1:8080
```

The server serves `web/dist` when a production frontend build exists. If `web/dist` is absent, it falls back to the source files under `web`.

## Frontend Development

Run the Go backend:

```sh
go run ./cmd/web-terminal
```

In another terminal, run Vite:

```sh
npm --prefix web run dev
```

Vite proxies `/api`, `/healthz`, and `/ws` to the local Go server.

## Production Build

```sh
npm --prefix web run build
go build -o web-terminal ./cmd/web-terminal
./web-terminal
```

Set `WEB_TERMINAL_ADDR` to use a different local bind address, and `WEB_TERMINAL_WEB_DIR` to serve a specific frontend directory.

## OpenAI Assistant

The assistant can use OpenAI from the local Go backend. The API key is read from your environment and is never sent to the browser.

Create a local `.env` file:

```sh
cp .env.example .env
```

Then edit `.env`:

```text
OPENAI_API_KEY=your_api_key_here
OPENAI_MODEL=gpt-5.4-mini
```

Run the app:

```sh
go run ./cmd/web-terminal
```

You can also export variables manually:

```sh
export OPENAI_API_KEY="your_api_key_here"
go run ./cmd/web-terminal
```

Optional model override:

```sh
export OPENAI_MODEL="gpt-5.4-mini"
```

If `OPENAI_API_KEY` is not set, or if the OpenAI request fails, the app falls back to the local rule-based suggestions.

Do not commit API keys. `.env` and `.env.*` files are ignored by Git, and `.env.example` is safe to commit as a template.

## Safety

This app controls a real shell. Treat every command as if you typed it in your normal terminal.

- Assistant suggestions are previews, not automatic actions.
- OpenAI-powered suggestions send the assistant request text to OpenAI. The API key stays on the local backend.
- Medium risk commands show warnings and require an explicit Run click.
- High risk commands require typing the exact command before running.
- Approved assistant commands are logged locally.
- The default audit log path is `~/.web-terminal/audit.log`.
- For tests or custom deployments, set `WEB_TERMINAL_AUDIT_LOG=/path/to/audit.log`.

## Test

```sh
go test ./...
npm --prefix web run build
```

## Troubleshooting

- Port already in use: run with `WEB_TERMINAL_ADDR=127.0.0.1:8090 go run ./cmd/web-terminal`.
- Blank or stale UI after frontend changes: run `npm --prefix web run build` again.
- Terminal does not connect: confirm the Go server is running and the browser is opened on the same local address.
- Assistant audit fails: check permissions for `~/.web-terminal` or set `WEB_TERMINAL_AUDIT_LOG` to a writable path.
- Shell is unexpected: the backend uses `$SHELL` when available and falls back to `/bin/sh`.

## Current Non-Goals

- Hosted remote terminal access.
- Multi-user authentication.
- Full natural language shell planning.
- Audit log viewer or retention management.
