# Agent 5: Safety and Audit

## Mission

Add safety controls around command execution, especially assistant-generated commands. The app should make dangerous actions visible, require confirmation, and keep a local audit trail.

Start after Agent 1 has created the skeleton.

## Ownership

Primary files and folders:

- `internal/safety/risk.go`
- `internal/safety/policy.go`
- `internal/audit/log.go`
- `internal/server/routes.go`
- `web/src/assistant.js`
- `web/src/styles.css`

Shared files that may need small edits:

- `internal/assistant/suggest.go`
- `web/src/api.js`

## Deliverables

- `POST /api/commands/risk` endpoint.
- Risk detection for shell commands.
- Confirmation requirements for risky commands.
- Local audit logging for approved assistant commands.
- UI warning states for medium and high risk suggestions.

## Risk Levels

### Low Risk

Examples:

- `pwd`
- `ls`
- `git status`
- `cat README.md`

### Medium Risk

Examples:

- Commands that overwrite files with `>`
- Commands that move files
- Package install commands
- Network commands such as `curl` or `wget`

### High Risk

Examples:

- `rm -rf`
- `sudo`
- `chmod -R`
- `chown -R`
- Disk formatting commands
- Commands that write to system directories
- Shell piping remote content into an interpreter, such as `curl ... | sh`

## Confirmation Policy

- Low risk: no extra confirmation.
- Medium risk: show warning and require explicit Run click.
- High risk: require a stronger confirmation step.
- Never allow the assistant to bypass confirmation.

## Audit Log

Store a local audit record for assistant-approved commands.

Suggested fields:

```json
{
  "timestamp": "2026-04-29T21:30:00Z",
  "cwd": "/Users/name/project",
  "command": "rm -rf dist",
  "risk": "high",
  "source": "assistant"
}
```

Suggested storage:

```text
~/.web-terminal/audit.log
```

## Done When

- Risk endpoint classifies common commands.
- Assistant UI displays risk warnings.
- High risk commands require explicit confirmation.
- Approved assistant commands are written to an audit log.
- Tests cover representative low, medium, and high risk commands.

