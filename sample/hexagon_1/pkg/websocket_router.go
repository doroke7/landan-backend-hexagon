package pkg

import "net/http"

// RouteGroup 模仿 gin 的 Group：把一組路由共用的路徑前綴集中管理，
// 呼叫端註冊時只要填寫相對路徑，實際掛上 mux 時會自動補上前綴。
type WebsocketRouter struct {
	prefix string
}

func NewWebsocketRouter(sPrefix string) *WebsocketRouter {
	return &WebsocketRouter{prefix: sPrefix}
}

func (oSelf *WebsocketRouter) HandleFunc(sPath string, fnHandler http.HandlerFunc) {
	http.HandleFunc(oSelf.prefix+sPath, fnHandler)
}
