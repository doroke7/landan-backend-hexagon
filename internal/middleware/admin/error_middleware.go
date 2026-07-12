package middleware_admin

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"example/pkg"

	"example/internal/bootstrap"
	"example/internal/utility"
)

type ErrorMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewErrorMiddleware(oAbstractMiddleware *AbstractMiddleware) *ErrorMiddleware {
	return &ErrorMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
// 放在中间件链最前面，先 Next() 让后续中间件/handler 执行，
// 执行完毕后统一检查是否有错误，有则返回统一格式的 JSON 响应。
func (oSelf *ErrorMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {
		oContext.Set("result", struct{}{}) // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("code", 0)            // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("message", "未知訊息")    // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("status", 200)        // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("authorization", "")  // 錯誤訊息給一個預設的，避免取 空 報錯

		defer func() {
			bHasError := false

			// 捕获 panic
			if oError := recover(); oError != nil {
				oByteStack := make([]byte, 4096)

				switch oErrorType := oError.(type) {
				case *pkg.DefaultError: // 需要用 *指標， 因為 controller 是用 指標
					fmt.Printf("[ERROR] %v", oErrorType)

					oSelf.Response.Set(oContext, 200, int(oErrorType.Code), oErrorType.Message, struct{}{}, "")

				default:
					iLen := runtime.Stack(oByteStack, false)
					fmt.Printf("[ERROR] %v\n%s\n", oError, oByteStack[:iLen])

					oSelf.Response.Set(oContext, 200, -4, "系統錯誤", struct{}{}, "")

				}

				bHasError = true
			}

			// 捕获 oContext.Error() 抛出的错误
			if len(oContext.Errors) > 0 {
				oLastErr := oContext.Errors.Last()

				byStack := make([]byte, 4096)
				iLen := runtime.Stack(byStack, false)
				fmt.Printf("[ERROR] %s\n%s\n", oLastErr.Error(), byStack[:iLen])

				oSelf.Response.Set(oContext, 200, -4, "系統錯誤", struct{}{}, "")

				bHasError = true
			}

			if !bHasError {
				return
			}

			mStatus, _ := oContext.Get("status")
			mCode, _ := oContext.Get("code")
			mResult, _ := oContext.Get("result")
			mMessage, _ := oContext.Get("message")

			sKey := utility.RandString(16)
			sIv := utility.RandString(16)

			sCode := fmt.Sprintf("%d", mCode)
			sMessage := mMessage.(string)
			iStatus := mStatus.(int)

			oKeys := map[string]interface{}{
				"key": sKey,
				"iv":  sIv,
			}
			sKeys, _ := utility.JsonEncode(oKeys)

			sTime := utility.Time[string](false)
			sResultJson, _ := utility.JsonEncode(mResult)

			sR := oSelf.aesHelper.Encrypt(sResultJson, sKey, sIv)
			sC := oSelf.aesHelper.Encrypt(sCode, sKey, sIv)
			sM := oSelf.aesHelper.Encrypt(sMessage, sKey, sIv)

			aStrings := []string{sKeys, sTime, sC, sM, sR, bootstrap.CONFIG.ADMIN.SIGNATURE.SALT}
			sHeaderSignature := utility.Md5(strings.Join(aStrings, ","))

			oContext.Writer.Header().Set("Authorization", "")
			oContext.Writer.Header().Set("Time", sTime)
			oContext.Writer.Header().Set("Signature", sHeaderSignature)

			oJson := gin.H{
				"c": sC,
				"m": sM,
				"r": sR,
			}
			if bootstrap.CONFIG.DEFAULT.DEBUG {
				oJson["code"] = mCode
				oJson["message"] = sMessage
				oJson["result"] = mResult
			}

			oContext.JSON(iStatus, oJson)
		}()

		oContext.Next()
	}
}
