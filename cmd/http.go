package cmd

import (
	"example/internal/bootstrap"
	"example/internal/container"
	"example/internal/register"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

var oHttpCommand = &cobra.Command{
	Use:   "http",
	Short: "啟動 Gin HTTP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitContainer()
		if err != nil {
			log.Fatal(err)
		}
		oGin := gin.Default()

		oEngine := register.HttpInit(oGin, oContainer)
		log.Fatal(oEngine.Run(":" + bootstrap.CONFIG.HTTP.PORT))
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oHttpCommand)
}
