package cmd

import (
	"context"
	"example/internal/container"
	"example/internal/register"

	"github.com/spf13/cobra"
)

var oSocketCommand = &cobra.Command{
	Use:   "socket",
	Short: "啟動 socket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitContainer()

		oClientRouter := register.ClientInit(oContainer)
		oClientRouter.Serve(context.Background())
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSocketCommand)
}
