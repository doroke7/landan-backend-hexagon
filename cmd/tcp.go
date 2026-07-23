package cmd

import (
	"context"
	bootstrap "example/bootstrap"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oTcpCommand = &cobra.Command{
	Use:   "tcp",
	Short: "啟動 TCP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitTcpContainer()
		if err != nil {
			log.Fatal(err)
		}

		oTcpRouter := register.TcpInit(oContainer)

		fmt.Println("TCP Server Start :" + bootstrap.CONFIG.SERVICES.TCP.PORT)

		if err := oTcpRouter.Serve(context.Background(), ":"+bootstrap.CONFIG.SERVICES.TCP.PORT); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oTcpCommand)
}
