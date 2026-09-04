package startday

import (
	"fmt"
	"os"
	"time"

	"dj/internal/directory"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "startday",
		Short: "Start a new journal entry for the day",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
	now := time.Now()
	baseDir := viper.GetString("directory")

	err := directory.CheckDirectory(baseDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePath := directory.GetTodaysFileName(baseDir)

	// TODO: Move to directory module
	_, err = os.Stat(filePath)
	if err == nil {
		fmt.Println("Journal entry already exists for today")
		os.Exit(1)
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	header := fmt.Sprintf("# %s\n\n", now.Format("January 02, 2006 (01/02/06)"))
	_, err = file.WriteString(header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Journal entry created for %s-%s-%s\n", now.Format("02"), now.Format("Jan"), now.Format("2006"))
}
