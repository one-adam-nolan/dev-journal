package show

import (
	"context"
	"fmt"
	"os"

	"dj/cli/internal/display"
	djv1 "dj/gen/go/dj/v1"
	"dj/internal/journal"

	"github.com/spf13/cobra"
)

// NewCommands returns the show command group.
func NewCommands(svc *journal.Service, highlighter display.Highlighter) *cobra.Command {
	todayCmd := &cobra.Command{
		Use:   "today",
		Short: "Prints activity for the current day",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printDay(cmd.Context(), svc, highlighter, &djv1.GetDayRequest{
				Selector: &djv1.GetDayRequest_Today{Today: true},
			})
		},
	}

	yesterdayCmd := &cobra.Command{
		Use:   "yesterday",
		Short: "Prints activity for the previous day",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printDay(cmd.Context(), svc, highlighter, &djv1.GetDayRequest{
				Selector: &djv1.GetDayRequest_Yesterday{Yesterday: true},
			})
		},
	}

	dateCmd := &cobra.Command{
		Use:     "date [MM/DD/YYYY]",
		Short:   "Displays entry details for the date specified",
		Example: "dj show date 06/01/2099",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return printDay(cmd.Context(), svc, highlighter, &djv1.GetDayRequest{
				Selector: &djv1.GetDayRequest_DateString{DateString: args[0]},
			})
		},
	}

	showRootCmd := &cobra.Command{
		Use:   "show",
		Short: "Displays content from a journal entry",
	}

	showRootCmd.AddCommand(todayCmd)
	showRootCmd.AddCommand(yesterdayCmd)
	showRootCmd.AddCommand(dateCmd)
	return showRootCmd
}

func printDay(ctx context.Context, svc *journal.Service, highlighter display.Highlighter, req *djv1.GetDayRequest) error {
	resp, err := svc.GetDay(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to open file, are you sure there are entries for that date?: %w", err)
	}

	return highlighter.HighlightMarkdown(os.Stdout, string(resp.GetDay().GetContent()))
}
