# AGENTS.md

Instructions for AI coding agents working on dev-journal.

## Project overview

dev-journal (`dj`) is a Go CLI for daily markdown work journals. Users create daily files, append timestamped entries and bullets, view past logs, and interact via a terminal UI.

- **Module:** `dj` (Go 1.20+)
- **CLI framework:** Cobra — commands live under `cmd/`, root wiring in `cmd/root.go`
- **Config:** Viper, TOML file at `~/.djconfig`
- **TUI:** tview/tcell in `internal/tui`
- **Layout:** `cmd/` (commands), `internal/` (app code), `pkg/` (shareable libs) — see [ADR-0001](docs/ADRs/0001-project-layout-internal-pkg-cmd.md)

## Setup commands

- Install binary: `make install` or `go install -v .`
- Build locally: `go build -o dj .`
- Run CLI: `go run . --help`
- Dev Container: open in VS Code with `.devcontainer/`

## Testing instructions

- Run all tests: `go test ./...`
- Run tests for a specific package: `go test ./cmd/config/...`
- Verify build after changes: `go build .`
- Smoke test CLI: `go run . --help`
- Config tests write to `~/.djconfig` — they clean up after themselves but require a writable home directory

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
- Keep `cmd/` handlers thin — business logic belongs in `internal/` or `pkg/`
- App-specific code goes in `internal/`, shareable code in `pkg/`
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

## Security considerations

- Config file path is `~/.djconfig` — do not log or expose its contents
- Journal directory is user-configured — never hardcode personal paths in source
- Do not commit `.DS_Store`, build artifacts, or local config files
