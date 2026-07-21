package cmd

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
)

var oHttpCommand = &cobra.Command{
	Use:   "http",
	Short: "啟動 Gin HTTP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitHttpContainer()
		if err != nil {
			log.Fatal(err)
		}
		oGin := gin.Default()

		oEngine := register.HttpInit(oGin, oContainer)
		log.Fatal(oEngine.Run(":" + bootstrap.CONFIG.SERVICES.HTTP.PORT))
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oHttpCommand)
}
