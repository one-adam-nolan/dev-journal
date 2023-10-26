package addlogic

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func AddBulletToFile(filepath string, msg string) error {
	//Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return errors.New("journal entry file not found")
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	timestamp := time.Now().Format("15:04")
	entry := fmt.Sprintf("\n* %s- %s", timestamp, msg)

	if _, err = f.WriteString(entry); err != nil {
		return err
	}

	return nil
}
