package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	container "example/internal/container"
)

var CommandCommand = &cobra.Command{
	Use:   "command",
	Short: "啟動 Command 命令",
	Run: func(oCmd *cobra.Command, args []string) {
		fmt.Println("cmd command")
	},
}

func init() {
	oContainer, err := container.InitCommandContainer()
	if err != nil {
		panic(err)
	}

	CommandCommand.AddCommand(oContainer.CommandAdminReourceAppUser.IncreaseBalance())
	oRootCommand.AddCommand(CommandCommand)
}
