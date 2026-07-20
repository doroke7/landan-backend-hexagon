package client

import (
	client "example/internal/client"
	helper "example/internal/helper"
)

// AbstractHandler 放 client 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（grpc / http / consumer）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*helper.AesHelper
	Client *client.Client
}

func NewAbstractHandler(oAesHelper *helper.AesHelper, oClient *client.Client) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
		Client:    oClient,
	}
}
