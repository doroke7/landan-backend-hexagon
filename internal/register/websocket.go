package register

import (
	"net/http"

	bootstrap "example/bootstrap"
	container "example/internal/container"
)

func WebsocketInit(oContainer *container.WebsocketContainer) *http.Server {
	return &http.Server{
		Addr: ":" + bootstrap.CONFIG.WEBSOCKET.PORT,
	}
}
