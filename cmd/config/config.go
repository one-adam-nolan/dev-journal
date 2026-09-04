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

	printCmd := &cobra.Command{
		Use:   "print",
		Short: "Print the content of the config file",
		RunE:  printConfig,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	configCmd.AddCommand(setDirCmd)
	configCmd.AddCommand(printCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.SetConfigType("toml")
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath(home)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a new one with default values
			viper.Set("directory", filepath.Join(home, "Documents", "Dev-Journal"))

			err := viper.SafeWriteConfigAs(filepath.Join(home, ".djconfig"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.ReadInConfig()
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

func printConfig(cmd *cobra.Command, args []string) error {
	// Read the config file
	configFile := viper.ConfigFileUsed()

	if configFile == "" {
		fmt.Println("No config file found.")
		return fmt.Errorf("please create a config file")
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read the config file: %s", err)
	}

	fmt.Println("Config file content:")
	fmt.Println("------------------------------")
	fmt.Println(string(content))
	fmt.Println("------------------------------")

	return nil
}
