package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var timetrackCli = &cobra.Command{
	Use:   "timetrack",
	Short: "Time tracking CLI",
	Long:  `TimeTrack, a CLI tool for keeping track of your time, with easy integration to Google Calendar.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to TimeTrack!")
	},
}

func main() {
	if err := timetrackCli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
