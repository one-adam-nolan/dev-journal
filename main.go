package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dj",
		Short: "Dev-Journal CLI app",
	}

	// config.InitConfig(rootCmd)

	startdayCmd := &cobra.Command{
		Use:   "startday",
		Short: "Start a new journal entry for the day",
		Run:   startDay,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	setdirCmd := &cobra.Command{
		Use:   "setdir [directory]",
		Short: "Set the directory for journal entries",
		Args:  cobra.ExactArgs(1),
		Run:   setDirectory,
	}

	addEntryCmd := &cobra.Command{
		Use:   "addentry [entry]",
		Short: "Append a new entry to today's markdown file",
		Args:  cobra.ExactArgs(1),
		RunE:  addEntry,
	}

	addBulletCmd := &cobra.Command{
		Use:   "addbullet [bullet]",
		Short: "Append a new bullet to the most recent entry",
		Args:  cobra.ExactArgs(1),
		RunE:  addBullet,
	}

	configCmd.AddCommand(setdirCmd)

	rootCmd.AddCommand(startdayCmd)
	rootCmd.AddCommand(addEntryCmd)
	rootCmd.AddCommand(addBulletCmd)
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

func setDirectory(cmd *cobra.Command, args []string) {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("Config file not found")

		// Create a default config file if it doesn't exist
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defaultConfigFile := filepath.Join(homeDir, ".djconfig")

		err = os.MkdirAll(filepath.Dir(defaultConfigFile), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, err := os.Create(defaultConfigFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		configContent := `# Dev-Journal config file
directory = "/path/to/journal"
`
		_, err = file.WriteString(configContent)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configFile = defaultConfigFile
	}

	// Set the directory value in the config file
	directory := args[0]
	viper.Set("directory", directory)

	// Write the updated config file
	err := viper.WriteConfigAs(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Reload the config file
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Directory set to: %s\n", directory)
}

func init() {
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("toml")

	viper.ReadInConfig()
}

func addEntry(cmd *cobra.Command, args []string) error {
	filepath := getTodaysFileName()

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Journal entry file not found")
		return nil
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	timestamp := time.Now().Format("15:04")
	entry := fmt.Sprintf("\n## %s\n\n### %s\n", timestamp, args[0])

	if _, err = f.WriteString(entry); err != nil {
		return err
	}

	fmt.Printf("Added entry to %s\n", filepath)

	return nil
}

func openFile() (*os.File, error) {
	filepath := getTodaysFileName()

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Journal entry file not found")
		return nil, err
	}

	return os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
}

func addBullet(cmd *cobra.Command, args []string) error {
	f, err := openFile()
	if err != nil {
		return err
	}

	defer f.Close()

	timestamp := time.Now().Format("15:04")
	entry := fmt.Sprintf("\n* %s- %s", timestamp, args[0])

	if _, err = f.WriteString(entry); err != nil {
		return err
	}

	fmt.Printf("Added entry %s\n", entry)

	return nil
}
