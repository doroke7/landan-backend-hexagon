package abstract

import "example/helper"

// AbstractHandler 放 input 層各個 adapter (http handler / consumer ...) 共用的依賴，
// 之後新增 adapter 只要 embed 這個 struct，就不用重複宣告欄位跟建構子參數。
type AbstractHandler struct {
	AesHelper *helper.AesHelper
}

func NewAbstractHandler(oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
	}
}
