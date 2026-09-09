# Multi-surface monorepo layout (CLI, backend, web)

* Status: accepted
* Date: 2026-09-09

## Context and Problem Statement

dev-journal is evolving from a CLI-only tool toward CLI, backend server, and web client surfaces. The codebase needed a layout that supports shared domain logic and protobuf contracts today while remaining easy to split into separate repositories later.

## Decision Drivers

* Single module and shared `internal/journal` use cases for CLI and backend
* Clear surface boundaries (CLI vs server vs web)
* Protobuf + Buf remain the cross-surface contract
* Avoid premature multi-repo overhead while staying split-ready

## Considered Options

* **Flat root layout (status quo)** — simple, but mixes CLI/server concerns and makes future extraction harder
* **Immediate polyrepo** — clean boundaries, but high overhead for a solo/small-team project
* **Split-ready monorepo** — shared contract + domain at root; surface-specific code under `cli/`, `backend/`, `web/`

## Decision Outcome

Chosen option: **split-ready monorepo**, because it preserves development velocity while making surface ownership explicit and enabling a future API repo split without redesign.

Layout:
- `proto/`, `gen/` — shared contract and generated code
- `internal/journal`, `internal/config`, `internal/logs` — shared domain
- `cli/` — `main.go`, Cobra commands, TUI, terminal display
- `backend/` — Connect HTTP server entrypoint and wiring
- `web/` — placeholder for future browser client (`gen/ts/` later)

### Consequences

* Good: CLI and backend reuse the same journal service and Connect handler
* Good: Web can join later via generated TS clients from the same proto module
* Good: `proto/` can be extracted to a standalone API repo with minimal churn
* Bad: More top-level directories to navigate
* Bad: Import paths change for CLI packages (`dj/cli/...`)

## Related

* Implementation plan: [multi-surface-monorepo.md](../plans/multi-surface-monorepo.md)
* Builds on [ADR-0003](0003-tiered-architecture-buf-di.md)
