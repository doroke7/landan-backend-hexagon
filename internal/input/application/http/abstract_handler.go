package input_application_http

import (
	pkg "example/pkg"

	client "example/internal/client"
	helper "example/internal/helper"
)

// AbstractHandler 放 http 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / consumer）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*pkg.Response
	*helper.AesHelper
	*helper.JwtHelper
	*client.ResourceClient
}

func NewAbstractHandler(oResponse *pkg.Response, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper, oResourceClient *client.ResourceClient) *AbstractHandler {
	return &AbstractHandler{
		AesHelper:      oAesHelper,
		JwtHelper:      oJwtHelper,
		Response:       oResponse,
		ResourceClient: oResourceClient,
	}
}
