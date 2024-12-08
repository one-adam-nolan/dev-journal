package startday

import (
	"fmt"
	"os"
	"time"

	"github.com/one-adam-nolan/dev-journal/pkg/directory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StartdayCmd = &cobra.Command{
	Use:   "startday",
	Short: "Start a new journal entry for the day",
	Run:   startDay,
}

func startDay(cmd *cobra.Command, args []string) {
	now := time.Now()

	err := directory.CheckDirectory(getBaseDirectory())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePath := directory.GetTodaysFileName(getBaseDirectory())

	//TODO: Move to directory module
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

func getBaseDirectory() string {
	return viper.GetString("directory")
}
