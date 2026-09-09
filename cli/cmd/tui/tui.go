package tui

import (
	tuiapp "dj/cli/internal/tui"

	"github.com/spf13/cobra"
)

// NewCommand returns the tui command.
func NewCommand(factory tuiapp.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Open todays file in a TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := factory.New()
			return app.Display()
		},
	}
}
