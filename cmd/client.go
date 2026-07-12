package cmd

import (
	"context"
	"example/internal/container"
	"example/internal/register"

	"github.com/spf13/cobra"
)

var oClientCommand = &cobra.Command{
	Use:   "client",
	Short: "啟動 Client 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitContainer()

		oClientRouter := register.ClientInit(oContainer)
		oClientRouter.Serve(context.Background())
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oClientCommand)
}
