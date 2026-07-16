package register

import (
	"example/internal/container"

	"github.com/gin-gonic/gin"
)

func adminMiddlewares(oContainer *container.HttpContainer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// ALL middleware
		oContainer.AdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.AdminErrorMiddleware.Handle(),
		oContainer.AdminSignatureMiddleware.Handle(),
		oContainer.AdminDecryptionMiddleware.Handle(),
		oContainer.AdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.AdminResponseMiddleware.Handle(),
		oContainer.AdminEncryptionMiddleware.Handle(),
	}
}

func HttpInit(oGin *gin.Engine, oContainer *container.HttpContainer) *gin.Engine {

	oAdmin := oGin.Group("/Admin").Use(adminMiddlewares(oContainer)...)
	{
		oAdmin.POST("/Authentication/Authenticator/SignIn", oContainer.HttpAdminAuthenticationAuthenticator.SignIn)

		oAdmin.POST("/Resource/User/Add", oContainer.HttpAdminResourceUser.AddUser)
		oAdmin.GET("/Resource/User/Show", oContainer.HttpAdminResourceUser.ShowUser)
	}

	return oGin
}
