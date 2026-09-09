# YYYY-MMM Folder Format and probable-palm-tree Migration

> ADR: [0002-journal-folder-naming-yyyy-mmm.md](../ADRs/0002-journal-folder-naming-yyyy-mmm.md)

## Problem

Month folders are currently named `Month-Year` (e.g. `March-2024`, `Jan-2023`, `feb-2022`), which sorts alphabetically by month name instead of chronologically. The dev-journal CLI uses full month names via Go's `January-2006` layout in `internal/directory/directory.go`.

The README already describes a year-first layout but uses numeric months (`YYYY-MM`); we align it to **`YYYY-MMM`** format (`2024-Mar`, `2026-Sep`).

## Target layout

```
journal-root/
├── 2022-Feb/
│   └── 10-Thursday.md
├── 2024-Mar/
│   └── 27-Wednesday.md
├── Landmarks/          # unchanged
└── OST/                # unchanged structure; nested tp/ month folders also renamed
```

Go format string: `2006-Jan` (produces `Jan`, `Feb`, `Mar`, etc.)

## Part 1: dev-journal code changes

### 1. Centralize folder naming in `internal/directory`

- Add exported `MonthFolderName(t time.Time) string` returning `t.Format("2006-Jan")`
- Update `getThisMonthsFolder` to call `MonthFolderName(time.Now())`
- Fix existing bug in `GetFileContentFromDate`: resolve folder from parsed date, not today

### 2. Add `dj migrate folders` command

- **Usage:** `dj migrate folders [--dry-run] [path]` (defaults to configured journal directory)
- Recursively walk the target directory
- Match directories named `{Month}-{YYYY}` (full or abbreviated month, case-insensitive)
- Rename each to `YYYY-MMM`
- Skip directories already matching `^\d{4}-[A-Z][a-z]{2}$`
- Skip hidden dirs
- Dry-run prints planned renames; live run applies them
- Abort if target folder already exists

### 3. TUI folder ordering

Sort month folders descending by folder name (`2026-Sep` before `2025-Oct`). Non-matching dirs remain at the end in mtime order.

### 4. Documentation

- Update README with `YYYY-MMM/<DD-DayName>.md` layout
- Add `dj migrate folders` to commands table

## Part 2: probable-palm-tree migration

**Scope:** `/Users/adam_nolan/Documents/Code/probable-palm-tree`

- 39 top-level month folders + 7 nested under `OST/Projects/tp/`
- ~455 journal markdown files; no merge conflicts
- Non-journal dirs left untouched: `Landmarks/`, `OST/` (parent), `tutorials/`

### Migration steps

1. Build updated `dj` binary
2. Point config at probable-palm-tree
3. Run `dj migrate folders --dry-run` and verify
4. Run `dj migrate folders` to apply renames
5. Remove legacy Python/shell scripts
6. Smoke test: `dj startday`, `dj show today`, `dj tui`

### Note on within-year sorting

`YYYY-MMM` fixes cross-year sort. Within a single year, abbreviated months do not sort chronologically alphabetically (`Apr` before `Jan`). `YYYY-MM` is the alternative if that becomes an issue.

## Testing

- `go test ./internal/directory/...`
- `go test ./...`
- Manual dry-run + live migration on probable-palm-tree
- Verify historical date: `dj show date 03/27/2024` → `2024-Mar/27-Wednesday.md`
