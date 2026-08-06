package middleware_admin

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
)

type RequestMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewRequestMiddleware(oAbstractMiddleware *AbstractMiddleware) *RequestMiddleware {
	return &RequestMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *RequestMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		if bootstrap.CONFIG.DEFAULT.DEBUG {
			// 直接讀 c.Request.URL.Query()，不要用 c.DefaultQuery——理由跟
			// SignatureMiddleware/DecryptionMiddleware 一樣：避免提前把 gin 的
			// queryCache 卡死，下面改寫 RawQuery 之後才不用反過來清快取。
			oUrlQuery := oContext.Request.URL.Query()
			sQuerySearch := oUrlQuery.Get("search")
			sQueryOption := oUrlQuery.Get("option")

			if sQueryOption != "" {
				oOption, _ := utility.JsonDecode[struct {
					Size  string `json:"size"`
					Page  string `json:"page"`
					AppId string `json:"app_id"`
				}](sQueryOption)
				oSelf.Flatten(oUrlQuery, "option", oOption)

			}

			if sQuerySearch != "" {
				oSearch, _ := utility.JsonDecode[map[string]interface{}](sQuerySearch)

				oSelf.Flatten(oUrlQuery, "search", oSearch)

			}

			oContext.Request.URL.RawQuery = oUrlQuery.Encode()
		}
		oContext.Next()

	}
}
