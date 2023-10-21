package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "mycommand",
	Short: "This is my command",
	Long:  `This command does amazing things.`,
	Run: commandHandler,
}

func init() {
	timetrackCli.addCommand(myCmd)
}

func commandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("ADD!")
}