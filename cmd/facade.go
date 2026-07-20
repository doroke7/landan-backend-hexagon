package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/internal/container"
	register "example/internal/register"
)

var oFacadeCommand = &cobra.Command{
	Use:   "facade",
	Short: "啟動 Facade 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitFacadeContainer()
		oFacadeServer := register.FacadeInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.FACADE.PORT)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(oFacadeServer.Serve(oListener))
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oFacadeCommand)
}
