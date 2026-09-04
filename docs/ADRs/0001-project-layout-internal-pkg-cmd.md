# Project layout (internal/pkg/cmd)

* Status: accepted
* Date: 2026-09-04

## Context and Problem Statement

The dev-journal codebase had grown with packages scattered at the repository root (`add/`, `config/`, `show/`, `directory/`, `tui/`, `logs/`). This made the project harder to navigate after time away and blurred the line between CLI commands, app-specific logic, and potentially shareable libraries.

## Decision Drivers

* Go community conventions for CLI projects (`cmd/`, `internal/`, `pkg/`)
* Cobra's recommended layout (`main.go` at root, commands under `cmd/`)
* Desire to keep the root directory clean
* Enforce import boundaries: app code in `internal/` cannot be imported by external modules

## Considered Options

* **Flat root layout (status quo)** — simple for a small project, but scales poorly and mixes concerns at the top level
* **`cmd/dj/main.go` entry point** — standard Go binary layout, but unnecessary for a single-binary CLI when Cobra defaults to root `main.go`
* **`internal/` + `pkg/` + `cmd/` with root `main.go`** — separates commands, private app code, and shareable libraries while matching Cobra convention

## Decision Outcome

Chosen option: **`internal/` + `pkg/` + `cmd/` with root `main.go`**, because it aligns with both Go and Cobra conventions, keeps the root clean, and makes package intent explicit.

Layout:
- `cmd/` — Cobra command definitions, one package per command group
- `internal/` — app-specific code (`directory`, `tui`, `logs`)
- `pkg/` — shareable libraries (`addlogic`, `controls`)
- `main.go` — thin entry point calling `cmd.Execute()`

### Consequences

* Good: Clear separation of CLI surface area, private implementation, and reusable packages
* Good: `internal/` is enforced by the Go compiler — external importers cannot depend on app internals
* Good: Modular `cmd/` packages preserve existing command group structure and tests
* Bad: More directories to navigate, though shallower at the root
* Bad: All import paths change — one-time migration cost

## Related

* Implementation plan: [restructure-go-layout.md](../plans/restructure-go-layout.md)
