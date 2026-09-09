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
| `dj migrate folders [--dry-run] [path]` | Rename legacy `Month-Year` folders to `YYYY-MMM` |

Configuration is stored at `~/.djconfig` (TOML). On first run, a default journal directory of `~/Documents/Dev-Journal` is created.

## Installing

```bash
make build-cli
sudo mv dj /usr/local/bin
```

Or use the Makefile:

```bash
make install
```

Build the backend server:

```bash
make build-server
./dj-server
```

> If you use `go install ./cli` it will put the binary in `$GOBIN`, which may make it so the autocompletion does not work properly.

## Codebase

Go module: `dj` (Go 1.20+). Built with [Cobra](https://cobra.dev/docs/) for CLI commands, [Viper](https://github.com/spf13/viper) for config, [Buf](https://github.com/bufbuild/buf) for protobuf schemas, [Connect](https://connectrpc.com/) for RPC contracts, [tview](https://github.com/rivo/tview) for the TUI, and [chroma](https://github.com/alecthomas/chroma) for markdown syntax highlighting.

```
dev-journal/
├── proto/dj/v1/         # shared protobuf contract (CLI, backend, web)
├── gen/go/dj/v1/        # generated Go types and Connect handlers
├── gen/ts/              # future web client output
├── cli/                 # dj CLI binary
│   ├── main.go
│   ├── cmd/             # Cobra commands
│   ├── internal/        # CLI-only app wiring, TUI, display
│   └── pkg/controls/    # tview helpers
├── backend/             # dj-server Connect HTTP server
│   ├── cmd/server/
│   └── internal/
├── web/                 # future browser client
├── internal/            # shared domain (journal, config, logs)
└── docs/
    ├── plans/
    └── ADRs/
```

Journal files are organized as `YYYY-MMM/<DD-DayName>.md` (e.g. `2026-Sep/04-Friday.md`) under the configured directory.

### Architecture

Shared domain logic lives in root `internal/journal`. CLI wiring is in `cli/internal/app`; backend wiring is in `backend/internal/app`. Both reuse the same protobuf contract in `proto/dj/v1/` and Connect handler in `internal/journal`.

### Key packages

- **`cli/cmd/`** — thin Cobra handlers
- **`cli/internal/app`** — CLI composition root
- **`cli/internal/tui`** — terminal UI
- **`backend/cmd/server`** — Connect HTTP server entrypoint
- **`internal/journal`** — shared use cases, filesystem store, Connect adapter
- **`internal/config`** — Viper config provider

See [ADR-0004](docs/ADRs/0004-multi-surface-monorepo.md) for the multi-surface layout, [ADR-0003](docs/ADRs/0003-tiered-architecture-buf-di.md) for tiered architecture and Buf, and [ADR-0001](docs/ADRs/0001-project-layout-internal-pkg-cmd.md) for original package layout.

## Development

Install Buf (for proto changes):

```bash
brew install bufbuild/buf/buf
```

```bash
make generate   # buf generate (also runs before test/build)
make test
make build      # builds dj and dj-server
go fmt ./...
go run ./cli --help
```

Code follows [Effective Go](https://go.dev/doc/effective_go) conventions. See [AGENTS.md](AGENTS.md) for agent-specific style rules.

VS Code launch configurations are in `.vscode/launch.json`. A Dev Container setup is available under `.devcontainer/`.

## Documentation

- [Effective Go](https://go.dev/doc/effective_go) — idiomatic Go conventions
- [Architectural Decision Records](docs/ADRs/README.md)
- [Restructure plan](docs/plans/restructure-go-layout.md)
- [AGENTS.md](AGENTS.md) — instructions for AI coding agents
