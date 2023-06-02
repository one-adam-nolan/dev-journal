package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetTodaysFileName(baseDir string) string {

	thisMonthsFolder := getThisMonthsFolder(baseDir)

	fileName := fmt.Sprintf("%s.md", time.Now().Format("02-Monday"))

	return filepath.Join(thisMonthsFolder, fileName)
}

func CheckDirectory(baseDir string) error {
	thisMonthsFolder := getThisMonthsFolder(baseDir)

	err := os.MkdirAll(thisMonthsFolder, os.ModePerm)

	return err
}

func GetTodaysFileContent(baseDirectory string) ([]byte, error) {
	filePath := GetTodaysFileName(baseDirectory)

	return os.ReadFile(filePath)
}

func getThisMonthsFolder(baseDirectory string) string {
	thisMonthsFolder := filepath.Join(baseDirectory, time.Now().Format("January-2006"))

	return thisMonthsFolder
}
