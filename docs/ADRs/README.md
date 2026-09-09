# Architectural Decision Records

This directory contains [Architectural Decision Records (ADRs)](https://adr.github.io/) for dev-journal. ADRs capture significant design decisions and their rationale.

## Format

We use **MADR minimal** (Markdown Architectural Decision Records).

### Template

```markdown
# {short title}

* Status: proposed | accepted | deprecated | superseded by [ADR-NNNN](NNNN-title.md)
* Date: YYYY-MM-DD

## Context and Problem Statement

{What is the issue? What forces are at play?}

## Decision Drivers

* {driver 1}
* {driver 2}

## Considered Options

* {option 1}
* {option 2}

## Decision Outcome

Chosen option: "{option}", because {justification}.

### Consequences

* Good: {positive consequence}
* Bad: {negative consequence or trade-off}
```

## Conventions

- Files: `NNNN-kebab-case-title.md` (zero-padded, e.g. `0001-project-layout-internal-pkg-cmd.md`)
- Status lifecycle: `proposed` → `accepted` (or `rejected`); later `deprecated` / `superseded by ADR-XXXX`
- ADRs are immutable once accepted — add a new ADR to change direction, don't rewrite history
- One decision per ADR; link related ADRs in Context or Consequences
- Implementation details belong in `docs/plans/`, not ADRs

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [0001](0001-project-layout-internal-pkg-cmd.md) | Project layout (internal/pkg/cmd) | accepted |
| [0002](0002-journal-folder-naming-yyyy-mmm.md) | Journal folder naming (YYYY-MMM) | accepted |
| [0003](0003-tiered-architecture-buf-di.md) | Tiered architecture, Buf, and dependency injection | accepted |
| [0004](0004-multi-surface-monorepo.md) | Multi-surface monorepo layout (CLI, backend, web) | accepted |
