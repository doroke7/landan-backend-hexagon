package cmd

import (
	"example/internal/bootstrap"
	"example/internal/container"
	"example/internal/register"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var oResourceCommand = &cobra.Command{
	Use:   "resource",
	Short: "啟動 Resource 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitResourceContainer()
		if err != nil {
			log.Fatal(err)
		}
		oResourceServer := register.ResourceInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.RESOURCE.PORT)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(oResourceServer.Serve(oListener))
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oResourceCommand)
}
