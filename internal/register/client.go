package register

import (
	pkg "example/pkg"

	container "example/container"
)

func ClientInit(oContainer *container.ClientContainer) *pkg.ClientRouter {
	oRouter := pkg.NewClientRouter()

	return oRouter
}
