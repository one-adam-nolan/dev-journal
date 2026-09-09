package tui

import (
	"dj/internal/journal"
)

// Displayable represents a TUI application that can be shown.
type Displayable interface {
	Display() error
}

// Factory creates TUI applications.
type Factory interface {
	New() Displayable
}

type factory struct {
	journal *journal.Service
}

// NewFactory returns a TUI factory backed by the journal service.
func NewFactory(svc *journal.Service) Factory {
	return &factory{journal: svc}
}

func (f *factory) New() Displayable {
	return &TextModalWithQandEscLowerBar{
		Journal: f.journal,
	}
}
