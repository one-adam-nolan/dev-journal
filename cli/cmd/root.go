package cmd

import (
	"dj/cli/cmd/add"
	"dj/cli/cmd/config"
	"dj/cli/cmd/migrate"
	"dj/cli/cmd/show"
	"dj/cli/cmd/startday"
	"dj/cli/cmd/tui"
	"dj/cli/internal/app"

	"github.com/spf13/cobra"
)

// Execute runs the dev-journal CLI with injected application dependencies.
func Execute(application *app.App) error {
	rootCmd := &cobra.Command{
		Use:   "dj",
		Short: "Dev-Journal CLI app",
	}

	rootCmd.AddCommand(config.NewCommands(application.Config))
	rootCmd.AddCommand(add.NewCommands(application.Journal))
	rootCmd.AddCommand(show.NewCommands(application.Journal, application.Display))
	rootCmd.AddCommand(migrate.NewCommands(application.Journal, application.Config))
	rootCmd.AddCommand(startday.NewCommand(application.Journal))
	rootCmd.AddCommand(tui.NewCommand(application.TUI))

	return rootCmd.Execute()
}
