package register

import (
	"example/pkg"
	"net/http"

	bootstrap "example/internal/bootstrap"
	container "example/internal/container"
)

func WebsocketInit(oContainer *container.WebsocketContainer) *http.Server {
	oGroup := pkg.NewWebsocketRouter("/websocket")
	oGroup.HandleFunc("/user/add", oContainer.WebsocketUser.AddUser)

	return &http.Server{
		Addr: ":" + bootstrap.CONFIG.WEBSOCKET.PORT,
	}
}
