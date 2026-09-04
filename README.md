# Dev Journal (DJ)

A terminal CLI for keeping a daily markdown dev journal — add timestamped entries, browse history, and view your log from anywhere.

## Overview

It started with tracking my work in markdown files a few years back. I'd create a new markdown file and start typing away each day. The formatting could have been more consistent, timestamps were only sometimes there, or I would get too lazy to type it out.

After a while, I made a few Python scripts to spice things up. But I wrote those with relative paths, meaning I had to execute them in a specific directory. Also, my "source code" was in the journal directory itself.

While attempting to add a Vim script to execute the Python scripts from any terminal session anywhere, I realized I needed to create a CLI application to solve this problem.

This project solves a few problems:

* I get to work with `Go`, which is new to me
* Log my work from anywhere in a terminal, including with Vim scripts
* Work with Dev Containers

## Commands

| Command | Description |
|---------|-------------|
| `dj startday` | Create today's journal markdown file |
| `dj add entry <text>` | Append a timestamped entry header |
| `dj add bullet <text>` | Append a bullet under the latest entry |
| `dj show today` | Print today's journal with syntax highlighting |
| `dj show yesterday` | Print yesterday's journal |
| `dj show date <MM/DD/YYYY>` | Print a specific date's journal |
| `dj tui` | Open today's journal in an interactive terminal UI |
| `dj config setdir <path>` | Set the journal directory |
| `dj config print` | Print the current config file |

Configuration is stored at `~/.djconfig` (TOML). On first run, a default journal directory of `~/Documents/Dev-Journal` is created.

## Installing

```bash
go build -o dj .
sudo mv dj /usr/local/bin
```

Or use the Makefile:

```bash
make install
```

> If you use `go install` it will put the binary in `$GOBIN`, which may make it so the autocompletion does not work properly.

## Codebase

Go module: `dj` (Go 1.20+). Built with [Cobra](https://cobra.dev/docs/) for CLI commands, [Viper](https://github.com/spf13/viper) for config, [tview](https://github.com/rivo/tview) for the TUI, and [chroma](https://github.com/alecthomas/chroma) for markdown syntax highlighting.

```
dev-journal/
├── main.go              # thin entry point → cmd.Execute()
├── cmd/                 # Cobra command definitions
│   ├── root.go          # root command wiring
│   ├── add/             # add entry / add bullet
│   ├── config/          # setdir, print
│   ├── show/            # today, yesterday, date
│   ├── startday/        # create today's file
│   └── tui/             # launch TUI
├── internal/            # app-specific, not importable externally
│   ├── directory/       # journal file paths, reads, folder listing
│   ├── tui/             # TUI implementation (modals, forms, history)
│   └── logs/            # colored logger (unused, reserved)
├── pkg/                 # shareable libraries
│   ├── addlogic/        # append entries/bullets to markdown files
│   └── controls/        # reusable tview button/tab helpers
└── docs/
    ├── plans/           # implementation plans
    └── ADRs/            # architectural decision records
```

Journal files are organized as `YYYY-MM/<DD-DayName>.md` (e.g. `2026-09/04-Thursday.md`) under the configured directory.

### Key packages

- **`cmd/`** — thin Cobra handlers; parse args, call into `internal/` or `pkg/`
- **`internal/directory`** — file naming, directory creation, content reads, folder sorting
- **`internal/tui`** — full-screen modal UI for viewing and editing today's journal
- **`pkg/addlogic`** — markdown append logic (`## HH:MM` entries, `* HH:MM-` bullets)
- **`pkg/controls`** — tview button styling and tab-focus navigation

See [ADR-0001](docs/ADRs/0001-project-layout-internal-pkg-cmd.md) for the rationale behind this layout.

## Development

```bash
go build .
go test ./...
go fmt ./...
go run . --help
```

Code follows [Effective Go](https://go.dev/doc/effective_go) conventions. See [AGENTS.md](AGENTS.md) for agent-specific style rules.

VS Code launch configurations are in `.vscode/launch.json`. A Dev Container setup is available under `.devcontainer/`.

## Documentation

- [Effective Go](https://go.dev/doc/effective_go) — idiomatic Go conventions
- [Architectural Decision Records](docs/ADRs/README.md)
- [Restructure plan](docs/plans/restructure-go-layout.md)
- [AGENTS.md](AGENTS.md) — instructions for AI coding agents
