package register

import (
	"example/internal/container"
	"example/pkg"
)

func ConsumerInit(oContainer *container.ConsumerContainer) *pkg.ConsumerRouter {
	oRouter := pkg.NewConsumerRouter(oContainer.Conn)

	return oRouter
}
