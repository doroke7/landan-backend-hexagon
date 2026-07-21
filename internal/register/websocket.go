package register

import (
	"net/http"

	bootstrap "example/bootstrap"
	container "example/container"
)

func WebsocketInit(oContainer *container.WebsocketContainer) *http.Server {
	return &http.Server{
		Addr: ":" + bootstrap.CONFIG.WEBSOCKET.PORT,
	}
}
