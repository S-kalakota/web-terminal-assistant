# Web Terminal Assistant

Web Terminal Assistant is a local-first web application that brings a terminal-style experience into the browser. The goal is to let a user interact with their actual computer from a website running locally on their machine.

The app is designed around a simple idea: the browser provides a friendlier interface, while a local Go backend handles the parts a normal website is not allowed to do, such as running shell commands or connecting to the filesystem.

## What It Does

- Serves a browser-based terminal interface from a local Go server.
- Runs locally at `http://127.0.0.1:8080`.
- Provides API routes for terminal communication, assistant suggestions, and command risk checks.
- Includes a JavaScript frontend shell with a terminal panel and assistant panel.
- Lays the foundation for a real PTY-backed terminal where commands like `pwd`, `ls`, `cd`, and `git status` can run on the user's computer.
- Includes the foundation for an English-to-command assistant that can suggest shell commands before the user approves them.

## Why It Runs Locally

A normal hosted website cannot safely access a user's hard drive or run terminal commands. Browsers block that on purpose.

This project solves that by running a local Go server on the user's computer. The browser talks to that local server, and the local server is responsible for terminal access.

By default, the app binds to:

```text
127.0.0.1:8080
```

## Tech Stack

- Go for the local backend server, routing, API contracts, and future terminal process management.
- JavaScript for the browser UI.
- HTML and CSS for the frontend layout and styling.
- Vite for frontend development and builds.
- npm for frontend dependency management.

## Project Structure

```text
cmd/web-terminal/        Go application entry point
internal/server/         HTTP server and route handlers
internal/terminal/       Shared terminal and API message structs
web/                     JavaScript frontend app
docs/                    Planning and implementation notes
```

## Run Locally

Install Go and Node.js, then from the project root run:

```sh
npm --prefix web install
go run ./cmd/web-terminal
```

Open:

```text
http://127.0.0.1:8080
```

## Build Frontend

```sh
npm --prefix web run build
```

## Current Status

This is an early project skeleton. The local server, static frontend, shared API message shapes, and placeholder assistant/risk routes are in place. The next major step is implementing the real PTY-backed terminal session so commands run through the user's actual shell.

## Safety

This project is intended to control a real shell on a real machine. Assistant-generated commands should always be previewed before running, and destructive commands should require extra confirmation.
