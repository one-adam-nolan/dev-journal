package cmds

import (
	"github.com/one-adam-nolan/dev-journal/cmds/add"
	"github.com/one-adam-nolan/dev-journal/cmds/config"
	"github.com/one-adam-nolan/dev-journal/cmds/show"
	"github.com/one-adam-nolan/dev-journal/cmds/startday"
	"github.com/one-adam-nolan/dev-journal/cmds/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dj",
	Short: "Dev-Journal CLI app",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	config.InitConfig(rootCmd)
	add.InitConfig(rootCmd)
	show.InitConfig(rootCmd)

	rootCmd.AddCommand(startday.StartdayCmd)
	rootCmd.AddCommand(tui.TuiCmd)
}
