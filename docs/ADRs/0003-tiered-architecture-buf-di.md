# Tiered architecture, Buf, and dependency injection

* Status: accepted
* Date: 2026-09-09

## Context and Problem Statement

The dev-journal CLI grew with global Viper config access, scattered file I/O across `cmd/`, `internal/directory`, and `pkg/addlogic`, and no service boundaries. Command handlers directly instantiated dependencies and some infrastructure code called `os.Exit`, making the codebase harder to test and extend.

## Decision Drivers

* Clear separation between CLI presentation, application use cases, domain contracts, and infrastructure
* Dependency injection from a single composition root (`main.go` / `internal/app`)
* Typed domain models and a versionable API contract via Protocol Buffers
* Connect-ready RPC handler for future remote journal access without rewriting use cases
* Testability via interfaces (`JournalStore`, `config.Provider`, `display.Highlighter`)

## Considered Options

* **Keep flat packages with global Viper** — minimal change, but does not address testability or API contracts
* **Go interfaces only, no proto** — improves DI but lacks a portable schema and RPC contract
* **Tiered architecture + proto/buf + DI composition root** — typed domain, Connect handler, injectable services

## Decision Outcome

Chosen option: **Tiered architecture + proto/buf + DI composition root**, because it establishes clear layer boundaries, wires dependencies in `main.go` via `internal/app`, defines journal operations in `proto/dj/v1/journal.proto`, and keeps markdown as the on-disk format while using generated types for service boundaries.

Layout:
- `proto/` + `buf.yaml` / `buf.gen.yaml` — schema source of truth
- `gen/go/` — generated Go protobuf and Connect code (committed)
- `internal/app` — composition root
- `internal/journal` — `JournalStore` interface, `FSJournalStore`, `Service`, `ConnectHandler`
- `internal/config` — `Provider` interface, Viper implementation
- `internal/display` — markdown highlighter interface
- `cmd/*` — thin handlers receiving injected dependencies

### Consequences

* Good: Commands are thin and testable; domain contract is explicit and lintable with Buf
* Good: Connect handler is ready for future HTTP server without changing CLI use cases
* Good: Infrastructure concerns (file I/O, config, highlighting) are behind interfaces
* Bad: More packages and generated code to maintain
* Bad: Proto changes require `buf generate` before build when modifying schemas

## Related

* Implementation plan: [tiered-architecture-buf.md](../plans/tiered-architecture-buf.md)
* Supersedes scattered logic previously in `internal/directory` and `pkg/addlogic`
