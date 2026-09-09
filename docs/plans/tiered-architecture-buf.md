# Tiered architecture with Buf and dependency injection

Related ADR: [ADR-0003](../ADRs/0003-tiered-architecture-buf-di.md)

## Goal

Refactor dev-journal into a tiered, dependency-injected architecture with protobuf domain types and a Connect-ready `JournalService` RPC contract.

## Layers

| Layer | Packages |
|-------|----------|
| Presentation | `main.go`, `cmd/*` |
| Application | `internal/app`, `internal/journal` (Service) |
| Domain | `proto/dj/v1/*` → `gen/go/dj/v1/` |
| Infrastructure | `internal/journal` (FSJournalStore), `internal/config`, `internal/display`, `internal/tui` |

## Dependency flow

```
main.go → app.New() → cmd.Execute(app)
cmd/* → journal.Service, config.Provider, display.Highlighter
journal.Service → JournalStore, config.Provider
FSJournalStore → filesystem + markdown helpers
```

## Proto contract

- `proto/dj/v1/config.proto` — `Config` message
- `proto/dj/v1/journal.proto` — journal messages + `JournalService` RPC

Generated with [Buf](https://github.com/bufbuild/buf):

```bash
make generate   # buf generate
```

## Wiring

All dependencies are composed in `internal/app/app.go` and passed into Cobra command constructors from `cmd/root.go`.

## Removed packages

- `internal/directory` — merged into `internal/journal/fs_store.go`
- `pkg/addlogic` — merged into `internal/journal/markdown.go`

## Future work

- Mount `journal.ConnectHandler` on HTTP for remote access
- BSR module publishing
