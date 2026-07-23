package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
)

var oFacadeCommand = &cobra.Command{
	Use:   "facade",
	Short: "啟動 Facade 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitFacadeContainer()
		if err != nil {
			log.Fatal(err)
		}
		oFacadeServer := register.FacadeInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.FACADE.PORT)
		if err != nil {
			log.Fatal(err)
		}

		// 收到中斷/終止訊號時主動 GracefulStop，讓 gRPC 停止 accept 新連線、
		// 關掉 listener，Serve() 才會正常返回並釋放 port，
		// 不是靠 process 被系統強制殺掉才釋放。
		oSignalChannel := make(chan os.Signal, 1)
		signal.Notify(oSignalChannel, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-oSignalChannel
			oFacadeServer.GracefulStop()
		}()

		if err := oFacadeServer.Serve(oListener); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oFacadeCommand)
}
