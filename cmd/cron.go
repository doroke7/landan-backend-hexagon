package cmd

import (
	"github.com/spf13/cobra"

	container "example/internal/container"
	register "example/internal/register"
)

var oCronCommand = &cobra.Command{
	Use:   "cron",
	Short: "啟動排程服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitCronContainer()

		oCron := register.CronInit(oContainer)
		oCron.Start()

		select {}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCronCommand)
}
