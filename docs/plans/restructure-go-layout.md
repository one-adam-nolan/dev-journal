# Restructure dev-journal to internal/ + pkg/ + cmd/

## Cobra convention

[Cobra's getting-started tutorial](https://cobra.dev/docs/tutorials/getting-started/) generates this layout:

```
my-cli/
  main.go          # thin entry point: calls cmd.Execute()
  cmd/
    root.go        # root command
    hello.go       # one file per command (cobra-cli add)
```

Keeping `main.go` at the root is the Cobra-recommended default for a single-binary CLI. Go's broader convention (`cmd/dj/main.go`) is also valid, but unnecessary here since we only ship one binary.

For larger command trees, Cobra docs describe a **modular pattern**: each command group in its own package under `cmd/`, exporting a constructor (the existing `InitConfig(rootCmd)` pattern). That fits this repo better than collapsing everything into one flat `cmd` package, because we already have separate command groups (`add`, `show`, `config`) with subcommands and tests.

## Target layout

```
dev-journal/
├── main.go                      # thin entry: wire commands, call Execute()
├── cmd/
│   ├── root.go                  # root cobra.Command + Execute()
│   ├── add/
│   │   ├── add.go
│   │   └── add_test.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── show/
│   │   └── show.go
│   ├── startday/
│   │   └── startday.go
│   └── tui/
│       └── tui.go               # cobra command only
├── internal/
│   ├── directory/
│   │   └── directory.go
│   ├── tui/                     # TUI implementation
│   │   ├── tui.go
│   │   ├── modals.go
│   │   ├── add_entry.go
│   │   ├── add_bullet.go
│   │   └── history_page.go
│   └── logs/
│       └── logger.go
├── pkg/
│   ├── addlogic/
│   └── controls/
├── go.mod
└── Makefile
```

## Package moves

| Previous path | New path | Rationale |
|---|---|---|
| `add/` | `cmd/add/` | Cobra command group |
| `config/` | `cmd/config/` | Cobra command group |
| `show/` | `cmd/show/` | Cobra command group |
| `main.go` `startDay` | `cmd/startday/` | Cobra command |
| `main.go` `displayTui` | `cmd/tui/` | Cobra command (imports `internal/tui`) |
| `directory/` | `internal/directory/` | App-specific, not shareable |
| `tui/` | `internal/tui/` | App-specific TUI implementation |
| `logs/` | `internal/logs/` | App-specific (unused, kept for future use) |
| `pkg/addlogic/` | unchanged | Shareable markdown append logic |
| `pkg/controls/` | unchanged | Shareable tview helpers |

## Import path updates

- `dj/add` → `dj/cmd/add`
- `dj/config` → `dj/cmd/config`
- `dj/show` → `dj/cmd/show`
- `dj/directory` → `dj/internal/directory`
- `dj/tui` → `dj/internal/tui` (in command handlers and TUI files)
- `dj/pkg/addlogic` and `dj/pkg/controls` stay the same

## main.go and cmd/root.go split

`main.go` is a thin entry point that calls `cmd.Execute()`.

`cmd/root.go` owns the root `cobra.Command`, wires all command groups, and handles errors.

`startDay` moved to `cmd/startday/startday.go`.
`displayTui` moved to `cmd/tui/tui.go` (imports `dj/internal/tui`).

## Cleanup

1. Removed duplicate empty `config` command registration from `main.go`
2. Removed dead `printToday` from `main.go`
3. Consolidated viper init — only `cmd/config`'s `cobra.OnInitialize(initConfig)` remains
4. Deleted empty `main_test.go` stub

## Supporting file updates

- `Makefile`: `go install -v .`
- `cmd/add/add_test.go`: import `dj/cmd/config`
- `.vscode/launch.json`: updated stale args

## Verification

```bash
go build .
go test ./...
go run . --help
go run . config print
```

## Related

- [ADR-0001: Project layout (internal/pkg/cmd)](../ADRs/0001-project-layout-internal-pkg-cmd.md)
