package app

import (
	"dj/internal/config"
	"dj/internal/journal"
	"dj/internal/logs"
)

// App holds wired dependencies for the backend server.
type App struct {
	Config  config.Provider
	Journal *journal.Service
	Connect *journal.ConnectHandler
	Logger  logs.Logger
}

// New wires and returns the backend dependency graph.
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
		Logger:  log,
	}, nil
}
