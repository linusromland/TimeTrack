package commands

import "github.com/urfave/cli/v2"

var AllCommands = []*cli.Command{
	AddCommand,
	SelectCalendarCommand,
	StartCommand,
	EndCommand,
	AbortCommand,
	ListCommand,
	InfoCommand,
}
