# Agent Workstream Index

This folder splits the web terminal CLI app into work packages that multiple people or coding agents can own without stepping on each other.

Start with Agent 1. Agent 1 creates the repo skeleton, shared contracts, and placeholder endpoints. Once that is merged, Agents 2 through 5 can work mostly in parallel. Agent 6 should run after the main implementation pieces are present.

## Workstreams

| Agent | Workstream | Main Ownership |
| --- | --- | --- |
| 1 | Skeleton and contracts | Repo layout, Go module, JS app shell, shared API types |
| 2 | Go terminal backend | Local server, WebSocket, PTY shell session |
| 3 | JavaScript terminal UI | Browser terminal experience and live WebSocket client |
| 4 | Assistant panel | English-to-command suggestion flow |
| 5 | Safety and audit | Risk detection, confirmations, local audit log |
| 6 | Packaging and QA | Build, run scripts, README, integration testing |

## Coordination Rules

- Keep ownership boundaries clear.
- Do not rewrite another agent's files unless the task explicitly requires it.
- If a shared contract must change, update `WEB_TERMINAL_CLI_APP_PLAN.md` and notify the other workstreams.
- Prefer small, testable slices over large rewrites.
- Keep the server local-only by default.
- Assistant-generated commands must never auto-run without user approval.

## Suggested Build Order

```text
1. Agent 1 creates the shared skeleton.
2. Agents 2, 3, 4, and 5 implement their components.
3. Agent 6 integrates, packages, and verifies the MVP.
```

