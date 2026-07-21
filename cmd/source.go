package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"

	pkg "example/pkg"

	bootstrap "example/bootstrap"
	container "example/internal/container"
	register "example/internal/register"
)

var oSourceCommand = &cobra.Command{
	Use:   "source",
	Short: "啟動 Rource 服務",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Logger(pkg.Default).Info("啟動 source 服務。 port: " + bootstrap.CONFIG.SERVICES.SOURCE.PORT)

		oContainer, err := container.InitResourceContainer()
		if err != nil {
			log.Fatal(err)
		}

		oResourceServer := register.SourceInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.SOURCE.PORT)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(oResourceServer.Serve(oListener))
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSourceCommand)
}
