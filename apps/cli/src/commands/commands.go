package commands

import (
	"TimeTrack-cli/src/app"

	"github.com/urfave/cli/v2"
)

func GetAllCommands(ctx *app.AppContext) []*cli.Command {
	return []*cli.Command{
		getAddTimeEntryCommand(ctx),
		getListTimeEntriesCommand(ctx),
		getLoginCommand(ctx),
		getRegisterCommand(ctx),
		getSettingsCommand(ctx),
	}
}
