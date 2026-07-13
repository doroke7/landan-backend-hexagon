package register

import (
	"example/internal/container"
	"example/pkg"
)

func ConsumerInit(oContainer *container.ConsumerContainer) *pkg.ConsumerRouter {
	oRouter := pkg.NewConsumerRouter(oContainer.ConsumerUser.Conn)
	oRouter.HandleFunc("User.Add", oContainer.ConsumerUser.AddUser)

	return oRouter
}
