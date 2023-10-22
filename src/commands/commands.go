package commands

import "github.com/urfave/cli/v2"

var AllCommands = []*cli.Command{
	AddCommand,
	SelectDatabaseCommand,
	StartCommand,
	EndCommand,
	AbortCommand,
	ListCommand,
}
