package cmd

import (
	"context"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oClientCommand = &cobra.Command{
	Use:   "client",
	Short: "啟動 Client 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitClientContainer()

		oClientRouter := register.ClientInit(oContainer)
		oClientRouter.Serve(context.Background())
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oClientCommand)
}
