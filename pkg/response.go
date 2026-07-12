package pkg

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (oSelf *Response) Set(oContext *gin.Context, iStatus int, iCode int, sMessage string, oResult any, sAuthorization string) {
	oContext.Set("code", iCode)
	oContext.Set("message", sMessage)
	oContext.Set("status", iStatus)
	oContext.Set("result", oResult)
	oContext.Set("authorization", sAuthorization)

}

func (oSelf *Response) SetWithNext(oContext *gin.Context, iStatus int, iCode int, sMessage string, oResult any, sAuthorization string) {
	oContext.Set("code", iCode)
	oContext.Set("message", sMessage)
	oContext.Set("status", iStatus)
	oContext.Set("result", oResult)
	oContext.Set("authorization", sAuthorization)

	oContext.Next()

}
