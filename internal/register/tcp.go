package register

import (
	pkg "example/pkg"

	container "example/container"
)

func TcpInit(oContainer *container.TcpContainer) *pkg.Tcp {
	oTcp := pkg.NewTcp()
	oTcp.HandleFunc("Admin.Authentication.Authenticator.SignIn", oContainer.TcpAdminAuthenticationSignIn.SignIn)

	return oTcp
}
