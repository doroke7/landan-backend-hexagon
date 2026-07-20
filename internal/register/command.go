package register

import (
	"github.com/spf13/cobra"

	container "example/internal/container"
)

func CommandInit(CommandCommand *cobra.Command, oContainer *container.CommandContainer) *cobra.Command {
	CommandCommand.AddCommand(oContainer.CommandAdminReourceAppUser.IncreaseBalance())

	return CommandCommand
}
