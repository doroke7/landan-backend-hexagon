package input_application

import (
	"github.com/gorilla/websocket"

	helper "example/internal/helper"
)

// AbstractHandler 放 websocket 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / http / consumer / cron）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*helper.AesHelper
	Upgrader websocket.Upgrader
}

func NewAbstractHandler(oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
		Upgrader:  websocket.Upgrader{},
	}
}
