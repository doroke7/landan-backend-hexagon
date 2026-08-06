package handler

import (
	helper "example/internal/helper"
	pkg "example/pkg"
)

// AbstractHandler 放 http 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / consumer）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*pkg.Response
	*helper.AesHelper
}

func NewAbstractHandler(oResponse *pkg.Response, oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
		Response:  oResponse,
	}
}
