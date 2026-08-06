package middleware_admin

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
)

type ResponseMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewResponseMiddleware(oAbstractMiddleware *AbstractMiddleware) *ResponseMiddleware {
	return &ResponseMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *ResponseMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

		mStatus, _ := oContext.Get("status")
		mCode, _ := oContext.Get("code")
		mMessage, _ := oContext.Get("message")
		mResult, _ := oContext.Get("result")
		mAuthorization, _ := oContext.Get("authorization")

		mC, _ := oContext.Get("c")
		mM, _ := oContext.Get("m")
		mR, _ := oContext.Get("r")
		mA, _ := oContext.Get("a")

		nCode := mCode.(int)
		sMessage := mMessage.(string)
		sAuthorization := mAuthorization.(string)

		sA := mA.(string)
		sC := mC.(string)
		sM := mM.(string)
		sR := mR.(string)

		iStatus := mStatus.(int)

		oJson := gin.H{
			"c": sC,
			"m": sM,
			"r": sR,
		}

		if bootstrap.CONFIG.DEFAULT.DEBUG {
			oJson["code"] = nCode
			oJson["message"] = sMessage
			oJson["result"] = mResult
			oContext.Writer.Header().Set("Authorization", sAuthorization)

		}
		oContext.Writer.Header().Set("A", sA)

		oContext.JSON(iStatus, oJson)
	}
}
