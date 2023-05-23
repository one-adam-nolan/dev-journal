package add

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(rootCmd *cobra.Command) {

	addEntryCmd := &cobra.Command{
		Use:   "entry [entry]",
		Short: "Append a new entry to today's markdown file",
		Args:  cobra.ExactArgs(1),
		RunE:  addEntry,
	}

	addBulletCmd := &cobra.Command{
		Use:   "bullet [bullet]",
		Short: "Append a new bullet to the most recent entry",
		Args:  cobra.ExactArgs(1),
		RunE:  addBullet,
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add an entry or a bullet point",
	}

	addCmd.AddCommand(addEntryCmd, addBulletCmd)

	rootCmd.AddCommand(addCmd)
}

func addEntry(cmd *cobra.Command, args []string) error {
	filepath := getTodaysFileName()

	//Check if file exists
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
	entry := fmt.Sprintf("\n\n## %s\n\n### %s", timestamp, args[0])

	if _, err = f.WriteString(entry); err != nil {
		return err
	}

	fmt.Printf("Added entry to %s\n", filepath)

	return nil
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
func openFile() (*os.File, error) {
	filepath := getTodaysFileName()

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Journal entry file not found")
		return nil, err
	}

	return os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
}
