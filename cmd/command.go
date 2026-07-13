package cmd

import (
	"fmt"

	container "example/internal/container"

	"github.com/spf13/cobra"
)

var CommandCommand = &cobra.Command{
	Use:   "command",
	Short: "啟動 Command 命令",
	Run: func(oCmd *cobra.Command, args []string) {
		fmt.Println("cmd command")
	},
}

func init() {
	var oContainer, _ = container.InitCommandContainer()

	CommandCommand.AddCommand(oContainer.UserHandler.AddUser())
	oRootCommand.AddCommand(CommandCommand)
}
