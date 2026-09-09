package config

import (
	"fmt"

	appconfig "dj/internal/config"

	"github.com/spf13/cobra"
)

// NewCommands returns the config command group.
func NewCommands(cfg appconfig.Provider) *cobra.Command {
	setDirCmd := &cobra.Command{
		Use:   "setdir [directory]",
		Short: "Set the directory for markdown files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetDir(cfg, args[0])
		},
	}

	printCmd := &cobra.Command{
		Use:   "print",
		Short: "Print the content of the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrintConfig(cfg)
		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	configCmd.AddCommand(setDirCmd)
	configCmd.AddCommand(printCmd)
	return configCmd
}

func runSetDir(cfg appconfig.Provider, directory string) error {
	if err := cfg.SetJournalDirectory(directory); err != nil {
		return err
	}

	fmt.Printf("Directory set to %s\n", directory)
	return nil
}

func runPrintConfig(cfg appconfig.Provider) error {
	content, err := cfg.RawConfig()
	if err != nil {
		return fmt.Errorf("failed to read the config file: %w", err)
	}

	fmt.Println("Config file content:")
	fmt.Println("------------------------------")
	fmt.Println(string(content))
	fmt.Println("------------------------------")
	return nil
}
