# Web Terminal Assistant

Web Terminal Assistant is a local-first web app concept for interacting with your real computer through a browser-based terminal. The long-term goal is a friendlier terminal experience: users can type normal shell commands, or describe what they want in English and review a suggested command before running it.

The app is intentionally designed as a local Go server plus a JavaScript browser UI. A regular hosted website cannot safely access a user's hard drive or run terminal commands, so this project runs on `127.0.0.1` and talks only to the machine where it is launched.

## What This Project Is Building

- A browser-based terminal connected to a local Go backend.
- A future PTY-backed shell session so commands like `pwd`, `cd`, `ls`, and `git status` run on the user's actual computer.
- An assistant panel that translates plain English into command suggestions.
- Safety checks for destructive commands before they run.
- Local audit logging for assistant-approved commands.
- A multi-agent implementation plan so different contributors can work on clear components.

## Current Status

This repository currently contains the Agent 1 skeleton from [docs/agents/agent-01-skeleton.md](docs/agents/agent-01-skeleton.md).

Implemented:

- Go module and local HTTP server.
- Static JavaScript frontend shell.
- Shared JSON message structs.
- `GET /healthz`.
- Placeholder `GET /ws/terminal`.
- Placeholder `POST /api/assistant/suggest`.
- Placeholder `POST /api/commands/risk`.
- Multi-agent planning docs.

Planned next:

- Real PTY-backed terminal session.
- `xterm.js` terminal UI.
- Complete assistant panel.
- Full command safety and audit logging.
- Production packaging.

## Why It Runs Locally

Browsers block direct access to the filesystem and shell for good reason. To make a real terminal-in-the-browser work, the user must run a local application that has permission to interact with the computer.

Default local address:

```text
http://127.0.0.1:8080
```

The server should stay bound to localhost unless explicit remote access and authentication are added later.

## Requirements

- Go 1.22 or newer.
- Node.js and npm for frontend development.

## Run The Skeleton

From the repository root:

```sh
go run ./cmd/web-terminal
```

Then open:

```text
http://127.0.0.1:8080
```

## Frontend Development

Install dependencies:

```sh
npm --prefix web install
```

Build the frontend:

```sh
npm --prefix web run build
```

## Environment

Optional settings:

```sh
WEB_TERMINAL_ADDR=127.0.0.1:8080
WEB_TERMINAL_WEB_DIR=web
```

## Current Routes

- `GET /`
- `GET /healthz`
- `GET /ws/terminal`
- `POST /api/assistant/suggest`
- `POST /api/commands/risk`

## Repository Plan

The build is split into agent-sized workstreams:

- Agent 1: Skeleton and contracts.
- Agent 2: Go terminal backend.
- Agent 3: JavaScript terminal UI.
- Agent 4: Assistant panel.
- Agent 5: Safety and audit.
- Agent 6: Packaging and QA.

See [docs/agents/README.md](docs/agents/README.md) for the full implementation split.

## Safety Note

This project is meant to control a real shell on a real machine. Assistant-generated commands must always be previewed and approved by the user before execution, and risky commands should require extra confirmation.
