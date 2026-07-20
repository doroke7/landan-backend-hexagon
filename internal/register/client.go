package register

import (
	"example/internal/container"
	"example/pkg"
)

func ClientInit(oContainer *container.ClientContainer) *pkg.ClientRouter {
	oRouter := pkg.NewClientRouter()

	return oRouter
}
