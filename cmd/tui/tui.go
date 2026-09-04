package tui

import (
	tuiapp "dj/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Open todays file in a TUI",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
	filePath := viper.GetString("directory")

	app := tuiapp.DisplayTodayModal(filePath)

	if err := app.Display(); err != nil {
		panic(err)
	}
}
