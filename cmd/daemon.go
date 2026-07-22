package cmd

import (
	"context"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oDaemonCommand = &cobra.Command{
	Use:   "daemon",
	Short: "啟動 daemon 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitDaemonContainer()

		oDaemonRouter := register.DaemonInit(oContainer)
		oDaemonRouter.Serve(context.Background())
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oDaemonCommand)
}
