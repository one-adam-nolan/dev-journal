package config

import (
	"testing"

	appconfig "dj/internal/config"

	"github.com/spf13/cobra"
)

func TestNewCommands(t *testing.T) {
	cmd := NewCommands(appconfig.NewViperProviderWithConfigDir(t.TempDir()))
	if cmd == nil {
		t.Fatal("expected config command")
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmd)

	if cmd.Use != "config" {
		t.Fatalf("expected config command, got %q", cmd.Use)
	}
}
