# Journal folder naming (YYYY-MMM)

* Status: accepted
* Date: 2026-09-04

## Context and Problem Statement

Journal month folders were named `Month-Year` (e.g. `March-2024`, `Jan-2023`, `feb-2022`). File explorers and directory listings sort these alphabetically by month name, not chronologically. Over time, inconsistent naming (full vs abbreviated month names, mixed case) made navigation worse.

## Decision Drivers

* Chronological sort in file explorers and the TUI history view
* Year-first ordering so recent months appear near the top when sorted descending
* A migration path for existing journals with legacy folder names
* Human-readable folder names (abbreviated month labels)

## Considered Options

* **Month-Year (status quo)** — familiar, but sorts by month name across years
* **YYYY-MM (numeric month)** — perfect chronological sort within and across years, less readable (`2024-03`)
* **YYYY-MMM (abbreviated month)** — year-first sort across years; readable (`2024-Mar`); chosen

## Decision Outcome

Chosen option: **YYYY-MMM**, using Go time layout `2006-Jan` (e.g. `2024-Mar`, `2026-Sep`) for all new month folders. Provide `dj migrate folders` to rename legacy `Month-Year` folders in existing journals.

### Consequences

* Good: Year-first alphabetical sort fixes cross-year navigation (`2023-Jan` before `2024-Mar`)
* Good: Consistent naming regardless of when folders were created
* Good: One-time migration command for existing journals
* Bad: Within a single year, abbreviated months do not sort chronologically alphabetically (`Apr` before `Jan`)
* Bad: Breaking change for existing journal directories until migration is run

## Related

* Implementation plan: [yyyy-mmm-folder-format.md](../plans/yyyy-mmm-folder-format.md)
