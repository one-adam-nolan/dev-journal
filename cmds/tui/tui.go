package tui

import (
	"github.com/one-adam-nolan/dev-journal/pkg/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var TuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open todays file in a TUI",
	Run:   displayTui,
}

func displayTui(cmd *cobra.Command, args []string) {
	filePath := getBaseDirectory()

	app := tui.DisplayTodayModal(filePath)

	// Run the application
	if err := app.Display(); err != nil {
		panic(err)
	}
}

func getBaseDirectory() string {
	return viper.GetString("directory")
}
