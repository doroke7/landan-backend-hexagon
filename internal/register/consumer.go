package register

import (
	pkg "example/pkg"

	container "example/internal/container"
)

func ConsumerInit(oContainer *container.ConsumerContainer) *pkg.ConsumerRouter {
	oRouter := pkg.NewConsumerRouter(oContainer.Conn)
	oRouter.HandleFunc("Admin.Resource.AppUser.IncreaseBalance", oContainer.ConsumerAdminResourceAppUser.IncreaseBalance)

	return oRouter
}
