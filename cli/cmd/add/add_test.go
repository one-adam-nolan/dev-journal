package add

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommands(t *testing.T) {
	cmd := NewCommands(nil)
	if cmd == nil {
		t.Fatal("expected add command")
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmd)

	if cmd.Use != "add" {
		t.Fatalf("expected add command, got %q", cmd.Use)
	}
}
