# Agent 6: Packaging and QA

## Mission

Integrate the workstreams into a runnable MVP, document how to use it, and verify that the core product flow works end to end.

Start after Agents 2 through 5 have working implementations.

## Ownership

Primary files and folders:

- `README.md`
- Build scripts or task runner files
- Frontend build configuration
- Go embed/static asset serving setup
- Integration test files
- Release notes or packaging docs

Shared files that may need small edits:

- `cmd/web-terminal/main.go`
- `internal/server/server.go`
- `web/package.json`

## Deliverables

- One-command local development flow.
- Production build flow that bundles the JS app with the Go server.
- README with setup, run, build, and safety notes.
- End-to-end verification checklist.
- Basic integration tests where practical.
- Known limitations section.

## Build Goals

Development:

```sh
go run ./cmd/web-terminal
```

Frontend development can use a separate dev server if needed, but the final app should be servable from the Go server.

Production:

```sh
npm --prefix web run build
go build -o web-terminal ./cmd/web-terminal
./web-terminal
```

## End-To-End QA Checklist

Verify:

- App starts on `127.0.0.1:8080`.
- Browser loads the UI.
- Terminal connects.
- Typing `pwd` shows the real current directory.
- Typing `cd ..` changes directory for future commands.
- Terminal resize works.
- Assistant can suggest `ls -laht` from English.
- Assistant suggestions do not run automatically.
- Approved suggestion runs in the terminal.
- High risk suggestion requires confirmation.
- Audit log records approved assistant commands.

## Documentation Requirements

README should include:

- What the app does.
- Why it must run locally.
- Install and run instructions.
- Development workflow.
- Safety warnings.
- Troubleshooting section.
- Current non-goals.

## Done When

- A new developer can clone, install, run, and use the MVP.
- Core terminal and assistant flow works end to end.
- Build commands are documented and tested.
- README clearly explains local-only security assumptions.

