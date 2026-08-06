package register

import (
	"example/pkg"
	"net/http"

	bootstrap "example/bootstrap"
	container "example/internal/container"
)

func WebsocketInit(oContainer *container.Container) *http.Server {
	oGroup := pkg.NewWebsocketRouter("/websocket")
	oGroup.HandleFunc("/user/add", oContainer.WebsocketUser.AddUser)

	return &http.Server{
		Addr: ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
	}
}
