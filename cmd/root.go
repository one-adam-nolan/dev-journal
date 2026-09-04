package cmd

import (
	"log"
	"os"

	"dj/cmd/add"
	"dj/cmd/config"
	"dj/cmd/show"
	"dj/cmd/startday"
	"dj/cmd/tui"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dj",
		Short: "Dev-Journal CLI app",
	}

	config.InitConfig(rootCmd)
	add.InitConfig(rootCmd)
	show.InitConfig(rootCmd)
	rootCmd.AddCommand(startday.NewCommand())
	rootCmd.AddCommand(tui.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
