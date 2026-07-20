package register

import (
	pkg "example/pkg"

	container "example/internal/container"
)

func ConsumerInit(oContainer *container.ConsumerContainer) *pkg.ConsumerRouter {
	oRouter := pkg.NewConsumerRouter(oContainer.Conn)
	oRouter.HandleFunc("AppUser.IncreaseBalance", oContainer.ConsumerAppUser.IncreaseBalance)

	return oRouter
}
