package register

import (
	pkg "example/pkg"

	container "example/internal/container"
)

func ClientInit(oContainer *container.ClientContainer) *pkg.ClientRouter {
	oRouter := pkg.NewClientRouter()

	return oRouter
}
