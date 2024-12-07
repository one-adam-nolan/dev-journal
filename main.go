package main

import (
	"log"
	"os"

	"github.com/one-adam-nolan/dev-journal/cmds"
	"github.com/spf13/viper"
)

func main() {

	err := cmds.Execute()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("toml")
	viper.ReadInConfig()
}
