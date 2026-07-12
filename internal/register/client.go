package register

import (
	"example/internal/container"
	"example/pkg"
)

func ClientInit(oContainer *container.Container) *pkg.ClientRouter {
	oRouter := pkg.NewClientRouter()
	oRouter.Handle(oContainer.ClientUser.AddUser)

	return oRouter
}
