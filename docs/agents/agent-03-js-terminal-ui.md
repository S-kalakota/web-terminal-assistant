# Agent 3: JavaScript Terminal UI

## Mission

Build the browser terminal experience. The user should see a polished terminal interface, type commands, and receive streamed output from the Go backend.

Start after Agent 1 has created the skeleton.

## Ownership

Primary files and folders:

- `web/src/main.js`
- `web/src/terminal.js`
- `web/src/api.js`
- `web/src/state.js`
- `web/src/styles.css`
- `web/index.html`
- `web/package.json`

Shared files that may need small edits:

- `README.md`

## Deliverables

- Browser UI with terminal area, assistant area placeholder, and top status bar.
- `xterm.js` integrated for terminal rendering.
- WebSocket connection to `GET /ws/terminal`.
- Terminal input sent to backend.
- Backend output written into terminal.
- Resize handling with terminal dimensions sent to backend.
- Connection status shown in the UI.
- Responsive layout for desktop and smaller screens.

## UI Requirements

Main layout:

```text
+---------------------------------------------------------------+
| Web Terminal       /Users/name/project            Connected   |
+--------------------------------------+------------------------+
|                                      | Assistant              |
| Terminal                             | Assistant panel        |
|                                      | placeholder or content |
+--------------------------------------+------------------------+
```

Design principles:

- Keep the terminal as the primary workspace.
- Keep UI chrome quiet and useful.
- Do not hide command output.
- Make connection status obvious.
- Leave room for Agent 4's assistant panel.

## Technical Notes

- Use `xterm.js`.
- Use a small API wrapper in `api.js` for WebSocket and HTTP calls.
- Keep terminal-specific behavior in `terminal.js`.
- Keep shared UI state in `state.js` if needed.
- Do not implement assistant logic here beyond a placeholder container. Agent 4 owns assistant behavior.

## Done When

- `npm install` and the selected frontend dev command work.
- The app connects to the Go WebSocket endpoint.
- Typing `pwd` in the browser terminal produces real output.
- Terminal resizes correctly.
- UI does not overlap at common desktop and mobile widths.

