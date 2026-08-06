package cmd

import (
	"example/internal/container"
	"example/internal/register"
	"log"

	"github.com/spf13/cobra"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitContainer()
		if err != nil {
			log.Fatal(err)
		}
		oWebsocketServer := register.WebsocketInit(oContainer)
		log.Fatal(oWebsocketServer.ListenAndServe())
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
