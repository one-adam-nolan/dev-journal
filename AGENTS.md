# AGENTS.md

Instructions for AI coding agents working on dev-journal.

## Project overview

dev-journal (`dj`) is a Go monorepo for daily markdown work journals — CLI today, backend server, and web client planned. Users create daily files, append timestamped entries and bullets, view past logs, and interact via a terminal UI.

- **Module:** `dj` (Go 1.20+)
- **Surfaces:** `cli/` (dj), `backend/` (dj-server), `web/` (future)
- **Shared domain:** `internal/journal`, `internal/config`, `internal/logs`
- **Contract:** protobuf in `proto/dj/v1/`, generated Go in `gen/go/`
- **CLI framework:** Cobra — commands live under `cli/cmd/`
- **Config:** Viper, TOML file at `~/.djconfig`
- **TUI:** tview/tcell in `cli/internal/tui`
- **Layout:** see [ADR-0004](docs/ADRs/0004-multi-surface-monorepo.md)

## Setup commands

- Install CLI: `make install` or `go install ./cli`
- Build CLI: `make build-cli` or `go build -o dj ./cli`
- Build server: `make build-server` or `go build -o dj-server ./backend/cmd/server`
- Run CLI: `go run ./cli --help`
- Run server: `go run ./backend/cmd/server`
- Dev Container: open in VS Code with `.devcontainer/`

## Testing instructions

- Run all tests: `make test` or `go test ./...`
- Run tests for a specific package: `go test ./cli/cmd/config/...`
- Verify build after changes: `make build`
- Smoke test CLI: `go run ./cli --help`
- Config tests use isolated temp directories via `NewViperProviderWithConfigDir` — they do not modify `~/.djconfig`

Add or update tests when changing command behavior or `pkg/` / `internal/` logic.

## Code style

Follow [Effective Go](https://go.dev/doc/effective_go) for idiomatic Go. Highlights for this project:

- Run `gofmt` (or `go fmt ./...`) — formatting is not negotiable
- Use mixedCaps / MixedCaps for names; short, concise identifiers; no `Get` prefix on getters
- Package names are lowercase, single-word, and match the directory name
- Return errors from functions (`RunE` in Cobra); avoid `panic` except for truly unrecoverable states
- Check errors — do not ignore return values; use `_` only when intentionally discarding
- Prefer small, focused interfaces; accept interfaces, return concrete types
- Write godoc-style comments on exported symbols: full sentences starting with the name
- Use `defer` for cleanup (file closes, etc.)
- Keep concurrency simple; this CLI is mostly synchronous — don't add goroutines without reason

Project-specific rules:

- Match existing conventions in the file you are editing
- Keep `cli/cmd/` handlers thin — business logic belongs in `internal/journal` or surface-specific `cli/internal/`
- Shared domain code goes in root `internal/`; CLI-only code in `cli/internal/`; backend wiring in `backend/internal/`
- New Cobra commands go under `cmd/` as their own package
- Do not add comments for obvious code; comment non-obvious business logic only
- Minimize scope — avoid unrelated changes in the same diff

## Project structure rules

- `main.go` stays a thin entry point calling `cmd.Execute()`
- Do not import `internal/` packages from outside this module (Go enforces this)
- Do not move shareable libraries out of `pkg/` without an ADR
- Significant architectural changes require a new ADR in `docs/ADRs/`

## Git and commit rules

- **Only developers may commit to git.** Agents must not run `git commit`, `git push`, or create pull requests unless explicitly asked by a developer.
- Agents may stage files or show diffs when helpful, but committing is a human action.
- Do not amend commits, force-push, or skip hooks unless a developer explicitly requests it.
- Do not commit secrets (`.env`, credentials, tokens) or generated binaries (`dj`).

## Documentation

- Human-facing overview: `README.md`
- Architectural decisions: `docs/ADRs/` (MADR minimal format)
- Implementation plans: `docs/plans/`
- Update `README.md` when adding user-facing commands or changing install steps
- Add an ADR when making structural or technology choices worth recording

## Planning and architectural decisions

When working on a non-trivial feature or structural change:

1. **Write a plan** in `docs/plans/` (kebab-case filename). If a plan was drafted in Cursor, copy the final version into `docs/plans/` before implementation begins.
2. **Write an ADR** in `docs/ADRs/` when the change involves a significant design choice (new conventions, layout changes, technology swaps, breaking behavior). Use the next sequential number and MADR minimal format. Update the index in `docs/ADRs/README.md`.
3. **Cross-link** the ADR and plan (ADR → plan in Related; plan → ADR at the top).
4. **Update `README.md`** when the change is user-facing (new commands, install steps, file layout).

## Security considerations

- Config file path is `~/.djconfig` — do not log or expose its contents
- Journal directory is user-configured — never hardcode personal paths in source
- Do not commit `.DS_Store`, build artifacts, or local config files
