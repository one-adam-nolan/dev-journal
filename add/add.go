package add

import (
	"dev-journal/directory"
	"dev-journal/pkg/addlogic"
	"fmt"

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
	filepath := directory.GetTodaysFileName(viper.GetString("directory"))

	err := addlogic.AddEntryToFile(filepath, args[0])
	if err != nil {
		fmt.Printf("Unable to add entry: %s\n", err)
	} else {
		fmt.Printf("Added entry: %s\n", filepath)
	}

	return nil
}

func addBullet(cmd *cobra.Command, args []string) error {
	filepath := directory.GetTodaysFileName(viper.GetString("directory"))

	err := addlogic.AddBulletToFile(filepath, args[0])
	if err != nil {
		fmt.Printf("Unable to add entry: %s\n", err)
	} else {
		fmt.Printf("Added bullet: %s\n", args[0])

	}

	return nil
}
