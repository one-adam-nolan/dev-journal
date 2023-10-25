package show

import (
	"dev-journal/directory"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(rootCmd *cobra.Command) {
	todayCmd := &cobra.Command{
		Use:   "today",
		Short: "Prints activity for the current day",
		Run:   printToday,
	}

	yesterdayCmd := &cobra.Command{
		Use:   "yesterday",
		Short: "Prints activity for the previous day",
		Run:   printYesterday,
	}

	dateCmd := &cobra.Command{
		Use:     "date",
		Short:   "Displays entry details for the date specified",
		Example: "dj show date 06/01/2099",
		Run:     printDate,
	}

	var showRootCmd = &cobra.Command{
		Use:   "show",
		Short: "Displays content from a journal entry",
	}

	showRootCmd.AddCommand(todayCmd)
	showRootCmd.AddCommand(yesterdayCmd)
	showRootCmd.AddCommand(dateCmd)
	rootCmd.AddCommand(showRootCmd)
}

func printToday(cmd *cobra.Command, args []string) {

	filePath := directory.GetTodaysFileName(getBaseDirectory())

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	_ = quick.Highlight(os.Stdout, string(content), "markdown", "terminal16m", "monokai")
}

func printYesterday(cmd *cobra.Command, args []string) {
	content, err := directory.GetFileContentFromDate(
		time.Now().
			AddDate(0, 0, -1).
			Format("01/02/2006"),
		getBaseDirectory())

	if err != nil {
		fmt.Printf("Unable to open file, are you sure there are entries for that date? \n")
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	_ = quick.Highlight(os.Stdout, string(content), "markdown", "terminal16m", "monokai")
}

func printDate(cmd *cobra.Command, args []string) {

	content, err := directory.GetFileContentFromDate(args[0], getBaseDirectory())
	if err != nil {
		fmt.Printf("Unable to open file, are you sure there are entries for that date? \n")
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	_ = quick.Highlight(os.Stdout, string(content), "markdown", "terminal16m", "monokai")
}

func getBaseDirectory() string {
	return viper.GetString("directory")
}
