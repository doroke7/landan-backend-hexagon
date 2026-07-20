package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	container "example/internal/container"
	register "example/internal/register"
)

var oConsumerCommand = &cobra.Command{
	Use:   "consumer",
	Short: "啟動 consumer 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitConsumerContainer()

		oConsumerRouter := register.ConsumerInit(oContainer)

		if err := oConsumerRouter.Serve(context.Background()); err != nil {
			log.Printf("consumer stopped: %v", err)
		}

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oConsumerCommand)
}
