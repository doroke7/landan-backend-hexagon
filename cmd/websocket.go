package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitWebsocketContainer()
		if err != nil {
			log.Fatal(err)
		}
		oWebsocketServer := register.WebsocketInit(oContainer)

		// 收到中斷/終止訊號時主動 Shutdown，讓 http.Server 停止 accept 新連線、
		// 關掉 listener，ListenAndServe() 才會正常返回並釋放 port，
		// 不是靠 process 被系統強制殺掉才釋放。
		oSignalChannel := make(chan os.Signal, 1)
		signal.Notify(oSignalChannel, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-oSignalChannel
			oWebsocketServer.Shutdown(context.Background())
		}()

		if err := oWebsocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
