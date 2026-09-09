package app

import (
	"dj/cli/internal/display"
	tuiapp "dj/cli/internal/tui"
	"dj/internal/config"
	"dj/internal/journal"
	"dj/internal/logs"
)

// App holds wired application dependencies for CLI commands.
type App struct {
	Config  config.Provider
	Journal *journal.Service
	Connect *journal.ConnectHandler
	Display display.Highlighter
	Logger  logs.Logger
	TUI     tuiapp.Factory
}

// New wires and returns the CLI dependency graph.
func New() (*App, error) {
	cfg := config.NewViperProvider()
	if err := cfg.Load(); err != nil {
		return nil, err
	}

	store := journal.NewFSStore(cfg)
	log := logs.NewColorfulLogger()
	journalSvc := journal.NewService(store, cfg, log)

	return &App{
		Config:  cfg,
		Journal: journalSvc,
		Connect: journal.NewConnectHandler(journalSvc),
		Display: display.NewChromaHighlighter(),
		Logger:  log,
		TUI:     tuiapp.NewFactory(journalSvc),
	}, nil
}
