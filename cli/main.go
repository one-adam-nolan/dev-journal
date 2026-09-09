package main

import (
	"log"
	"os"

	"dj/cli/cmd"
	"dj/cli/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(application); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
