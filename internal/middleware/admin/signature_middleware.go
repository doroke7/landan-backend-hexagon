package middleware_admin

import (
	"example/internal/bootstrap"
	"example/internal/utility"
	"strings"

	types "example/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"example/pkg"
)

type SignatureMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewSignatureMiddleware(oAbstractMiddleware *AbstractMiddleware) *SignatureMiddleware {
	return &SignatureMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *SignatureMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		sVer := oContext.GetHeader("Ver")
		sVersion := oContext.GetHeader("Version")
		sK := oContext.GetHeader("K")
		sTime := oContext.GetHeader("Time")
		sHeaderSignature := oContext.GetHeader("Signature")
		sS := oContext.DefaultQuery("s", "")
		sO := oContext.DefaultQuery("o", "")
		sP := oContext.PostForm("p")

		var oRequestPayload types.RequestPayload

		// 2. 關鍵：使用 c.ShouldBind 代替 c.ShouldBindJSON！
		// Gin 會自動根據 Content-Type 去選用 JSON 解析器或 Form 解析器
		if err := oContext.ShouldBindBodyWith(&oRequestPayload, binding.JSON); err != nil {
			oContext.Abort()
			_ = oContext.Error(pkg.NewDefaultError("請求格式錯誤1", -1, 400))

			return
		}
		sP = oRequestPayload.P

		// NOTE: 不要把 未加密的 search, option, param, 都加下去簽名，多次一舉
		aStrings := []string{sVer, sVersion, sK, sTime, sS, sO, sP, bootstrap.CONFIG.ADMIN.SIGNATURE.SALT}

		sStrings := strings.Join(aStrings, "|")
		sMd5Signature := utility.Md5(sStrings)

		if bootstrap.CONFIG.ADMIN.SIGNATURE.STATUS == true {
			if sMd5Signature != sHeaderSignature {
				oContext.Abort()
				_ = oContext.Error(pkg.NewDefaultError("簽名失敗", -3, 406))

				return
			}
		}

		oContext.Next()

	}
}
