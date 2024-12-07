package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Test_initConfig_creates_config(t *testing.T) {
	configWithPath := deleteConfigFile()

	initConfig()

	if _, err := os.Stat(configWithPath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Function did not create the default config file!")
	} else {
		t.Logf("Default config created successfully")
	}

	directory := viper.GetString("directory")

	if directory != "/root/Documents/Dev-Journal" {
		t.Fatalf("Incorrect directory in config file")
	} else {
		t.Logf("Directory set properly at: %s", directory)
	}
}

func Test_setdir_updates_config_file(t *testing.T) {
	deleteConfigFile()
	initConfig()

	directory := "/my/new/directory"

	cmd := &cobra.Command{}

	if err := setDir(cmd, []string{directory}); err != nil {
		t.Fatalf("Fail to set the directory, \n %s", err)
	}

	updatedValue := viper.GetString("directory")

	if updatedValue != directory {
		t.Fatalf("The new directorty was not set.")
	} else {
		t.Logf("Directory Updated Successfully: %s", directory)
	}

	deleteConfigFile()
	initConfig()
	t.Logf("Reset config file")

}

func deleteConfigFile() string {
	home, _ := os.UserHomeDir()
	configWithPath := filepath.Join(home, ".djconfig")
	os.Remove(configWithPath)
	return configWithPath
}
