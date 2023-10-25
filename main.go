package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dev-journal/add"
	"dev-journal/config"
	"dev-journal/directory"
	"dev-journal/show"
	"dev-journal/tui"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dj",
		Short: "Dev-Journal CLI app",
	}

	config.InitConfig(rootCmd)
	add.InitConfig(rootCmd)
	show.InitConfig(rootCmd)

	startdayCmd := &cobra.Command{
		Use:   "startday",
		Short: "Start a new journal entry for the day",
		Run:   startDay,
	}

	tuiCmd := &cobra.Command{
		Use:   "tui",
		Short: "Open todays file in a TUI",
		Run:   displayTui,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	rootCmd.AddCommand(startdayCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(configCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}

func getBaseDirectory() string {
	return viper.GetString("directory")
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

func printToday(cmd *cobra.Command, args []string) {

	filePath := directory.GetTodaysFileName(getBaseDirectory())

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	_ = quick.Highlight(os.Stdout, string(content), "markdown", "terminal16m", "monokai")
}

func displayTui(cmd *cobra.Command, args []string) {
	filePath := getBaseDirectory()

	app := tui.DisplayTodayModal(filePath)

	// Run the application
	if err := app.Display(); err != nil {
		panic(err)
	}
}

func init() {
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("toml")

	viper.ReadInConfig()
}
