package startday

import (
	"context"
	"fmt"
	"time"

	"dj/internal/journal"

	"github.com/spf13/cobra"
)

// NewCommand returns the startday command.
func NewCommand(svc *journal.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "startday",
		Short: "Start a new journal entry for the day",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), svc)
		},
	}
}

func run(ctx context.Context, svc *journal.Service) error {
	resp, err := svc.StartDay(ctx)
	if err != nil {
		if err == journal.ErrDayAlreadyExists {
			return fmt.Errorf("journal entry already exists for today")
		}
		return err
	}

	if !resp.GetCreated() {
		return fmt.Errorf("journal entry already exists for today")
	}

	now := time.Now()
	fmt.Printf("Journal entry created for %s-%s-%s\n", now.Format("02"), now.Format("Jan"), now.Format("2006"))
	return nil
}
