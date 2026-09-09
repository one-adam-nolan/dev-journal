package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dj/internal/config"
	"dj/internal/journal"

	"github.com/spf13/cobra"
)

// NewCommands returns the migrate command group.
func NewCommands(svc *journal.Service, cfg config.Provider) *cobra.Command {
	var dryRun bool

	foldersCmd := &cobra.Command{
		Use:   "folders [path]",
		Short: "Rename legacy Month-Year folders to YYYY-MMM",
		Long:  "Recursively renames journal month folders from legacy Month-Year format (e.g. March-2024) to YYYY-MMM (e.g. 2024-Mar). Defaults to the configured journal directory.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cfg.JournalDirectory()
			if len(args) == 1 {
				root = args[0]
			}

			root, err := filepath.Abs(root)
			if err != nil {
				return err
			}

			if _, err := os.Stat(root); err != nil {
				return fmt.Errorf("journal directory not found: %s", root)
			}

			return runMigrateFolders(cmd.Context(), svc, root, dryRun)
		},
	}

	foldersCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print planned renames without applying them")

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate journal data to new formats",
	}

	migrateCmd.AddCommand(foldersCmd)
	return migrateCmd
}

func runMigrateFolders(ctx context.Context, svc *journal.Service, root string, dryRun bool) error {
	resp, err := svc.MigrateFolders(ctx, root, dryRun)
	if err != nil {
		return err
	}

	renames := resp.GetRenames()
	if len(renames) == 0 {
		fmt.Println("No legacy month folders found.")
		return nil
	}

	for _, rename := range renames {
		fmt.Printf("%s -> %s\n", rename.GetOldPath(), rename.GetNewPath())
	}

	if dryRun {
		fmt.Printf("Dry run: %d folder(s) would be renamed.\n", len(renames))
	} else {
		fmt.Printf("Renamed %d folder(s).\n", len(renames))
	}

	return nil
}
