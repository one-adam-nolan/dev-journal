package add

import (
	"context"
	"fmt"

	"dj/internal/journal"

	"github.com/spf13/cobra"
)

// NewCommands returns the add command group.
func NewCommands(svc *journal.Service) *cobra.Command {
	addEntryCmd := &cobra.Command{
		Use:     "entry [entry]",
		Aliases: []string{"e"},
		Short:   "Append a new entry to today's markdown file",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddEntry(cmd.Context(), svc, args[0])
		},
	}

	addBulletCmd := &cobra.Command{
		Use:     "bullet [bullet]",
		Aliases: []string{"b"},
		Short:   "Append a new bullet to the most recent entry",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddBullet(cmd.Context(), svc, args[0])
		},
	}

	addCmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add an entry or a bullet point",
	}

	addCmd.AddCommand(addEntryCmd, addBulletCmd)
	return addCmd
}

func runAddEntry(ctx context.Context, svc *journal.Service, text string) error {
	resp, err := svc.AddEntry(ctx, text)
	if err != nil {
		return fmt.Errorf("unable to add entry: %w", err)
	}

	fmt.Printf("Added entry: %s\n", resp.GetFilePath())
	return nil
}

func runAddBullet(ctx context.Context, svc *journal.Service, text string) error {
	_, err := svc.AddBullet(ctx, text)
	if err != nil {
		return fmt.Errorf("unable to add bullet: %w", err)
	}

	fmt.Printf("Added bullet: %s\n", text)
	return nil
}
