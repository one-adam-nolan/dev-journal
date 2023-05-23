package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dev-journal/add"
	"dev-journal/config"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dj",
		Short: "Dev-Journal CLI app",
	}

	config.InitConfig(rootCmd)
	add.InitConfig(rootCmd)

	startdayCmd := &cobra.Command{
		Use:   "startday",
		Short: "Start a new journal entry for the day",
		Run:   startDay,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	rootCmd.AddCommand(startdayCmd)
	rootCmd.AddCommand(configCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// extract helper
func getTodaysFileName() string {

	thisMonthsFolder := getThisMonthsFolder()

	fileName := fmt.Sprintf("%s.md", time.Now().Format("02-Monday"))

	return filepath.Join(thisMonthsFolder, fileName)
}

// extract helper
func getThisMonthsFolder() string {
	folderPath := viper.GetString("directory")

	thisMonthsFolder := filepath.Join(folderPath, time.Now().Format("January-2006"))

	return thisMonthsFolder
}

// extract helper
func checkDirectory() error {
	thisMonthsFolder := getThisMonthsFolder()

	err := os.MkdirAll(thisMonthsFolder, os.ModePerm)

	return err
}

func startDay(cmd *cobra.Command, args []string) {
	now := time.Now()

	err := checkDirectory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePath := getTodaysFileName()

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

func init() {
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("toml")

	viper.ReadInConfig()
}
