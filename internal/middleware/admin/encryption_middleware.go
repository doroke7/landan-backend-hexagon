package middleware_admin

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
)

type EncryptionMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewEncryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *EncryptionMiddleware {
	return &EncryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *EncryptionMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

		mKey, _ := oContext.Get("key")
		mIv, _ := oContext.Get("iv")

		mCode, _ := oContext.Get("code")
		mResult, _ := oContext.Get("result")
		mMessage, _ := oContext.Get("message")
		mAuthorization, _ := oContext.Get("authorization")

		sCode := fmt.Sprintf("%d", mCode)
		sKey := fmt.Sprintf("%s", mKey)
		sIv := fmt.Sprintf("%s", mIv)
		sAuthorization := fmt.Sprintf("%s", mAuthorization)

		sMessage := mMessage.(string)

		sTime := utility.Time[string](false)
		sResult, _ := utility.JsonEncode(mResult)

		sR := oSelf.aesHelper.Encrypt(sResult, sKey, sIv)
		sC := oSelf.aesHelper.Encrypt(sCode, sKey, sIv)
		sM := oSelf.aesHelper.Encrypt(sMessage, sKey, sIv)
		sA := oSelf.aesHelper.Encrypt(sAuthorization, bootstrap.CONFIG.ADMIN.JWT.KEY, bootstrap.CONFIG.ADMIN.JWT.IV)

		// NOTE: 不要把 未加密的 code, message, result 都加下去簽名，多次一舉
		aStrings := []string{sTime, sC, sM, sR, bootstrap.CONFIG.ADMIN.SIGNATURE.SALT}

		sString := strings.Join(aStrings, ",")

		sHeaderSignature := utility.Md5(sString)

		oContext.Writer.Header().Set("Time", sTime)
		oContext.Writer.Header().Set("Signature", sHeaderSignature)

		oContext.Set("a", sA)
		oContext.Set("c", sC)
		oContext.Set("m", sM)
		oContext.Set("r", sR)

	}
}
