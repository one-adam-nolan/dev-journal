package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(rootCmd *cobra.Command) {
	cobra.OnInitialize(initConfig)

	setDirCmd := &cobra.Command{
		Use:   "setdir [directory]",
		Short: "Set the directory for markdown files",
		Args:  cobra.ExactArgs(1),
		RunE:  setDir,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	configCmd.AddCommand(setDirCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a new one with default values
			viper.Set("directory", filepath.Join("$HOME", "Documents", "Dev-Journal"))
			if err := viper.WriteConfigAs(filepath.Join("$HOME", ".djconfig")); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Created new config file at ~/.djconfig")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func setDir(cmd *cobra.Command, args []string) error {
	directory := args[0]

	viper.Set("directory", directory)

	configFile := viper.ConfigFileUsed()

	if err := viper.WriteConfigAs(configFile); err != nil {
		return err
	}

	fmt.Printf("Directory set to %s\n", directory)

	return nil
}
