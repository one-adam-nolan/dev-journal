# Multi-surface monorepo layout

Related ADR: [ADR-0004](../ADRs/0004-multi-surface-monorepo.md)

## Goal

Organize dev-journal as a split-ready monorepo for CLI, backend server, and future web client while keeping protobuf as the shared contract.

## Layout

```
dev-journal/
├── proto/dj/v1/           # shared API contract
├── gen/go/                # generated Go types (CLI + backend)
├── gen/ts/                # future web client output
├── cli/                   # dj binary
├── backend/               # dj-server binary
├── web/                   # future web app
└── internal/              # shared domain (journal, config, logs)
```

## Boundaries

| Path | Owns |
|------|------|
| `proto/`, `gen/` | Contract and generated code for all surfaces |
| `internal/journal`, `internal/config`, `internal/logs` | Shared domain and infrastructure |
| `cli/` | Cobra commands, TUI, terminal display |
| `backend/` | Connect HTTP server wiring |
| `web/` | Browser client (future) |

## Future repo split

If repos split later, extract `proto/` (+ buf config + gen) into `dev-journal-api` and pin versions from CLI, backend, and web.
