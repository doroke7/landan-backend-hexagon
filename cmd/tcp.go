package cmd

import (
	bootstrap "example/bootstrap"
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var oTcpCommand = &cobra.Command{
	Use:   "tcp",
	Short: "啟動 TCP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ln, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.TCP.PORT)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		fmt.Println("TCP Server Start")

		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}

			go func(conn net.Conn) {
				defer conn.Close()

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					fmt.Println("Receive:", string(buf[:n]))

					_, err = conn.Write([]byte("OK"))
					if err != nil {
						return
					}
				}
			}(conn)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oTcpCommand)
}
