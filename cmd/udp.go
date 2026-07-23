package cmd

import (
	"github.com/spf13/cobra"
)

var oUdpCommand = &cobra.Command{
	Use:   "udp",
	Short: "啟動 UDP 服務",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oUdpCommand)
}
