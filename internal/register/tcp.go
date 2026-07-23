package register

import (
	pkg "example/pkg"

	container "example/container"
)

func TcpInit(oContainer *container.TcpContainer) *pkg.TcpRouter {
	oRouter := pkg.NewTcpRouter()
	oRouter.HandleFunc("Admin.Authentication.Authenticator.SignIn", oContainer.TcpAdminAuthenticationSignIn.SignIn)

	return oRouter
}
