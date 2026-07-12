package middleware_admin

import (
	"example/internal/bootstrap"
	"example/internal/utility"
	"net/url"
	"reflect"
	"unsafe"

	"github.com/gin-gonic/gin"
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
			sQuerySearch := oContext.DefaultQuery("search", "")
			sQueryOption := oContext.DefaultQuery("option", "")
			oUrlQuery := oContext.Request.URL.Query()

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

			sEncoded := oUrlQuery.Encode()
			oContext.Request.URL.RawQuery = sEncoded

			/*
			   清理 Gin 的 queryCache，也就是強制清理 context 中的 get 參數緩存
			*/
			oContextValue := reflect.ValueOf(oContext).Elem()
			oQueryCacheField := oContextValue.FieldByName("queryCache")
			pQueryCacheField := oQueryCacheField.UnsafeAddr()
			ptrQueryCacheField := (*url.Values)(unsafe.Pointer(pQueryCacheField))
			*ptrQueryCacheField = nil
		}
		oContext.Next()

	}
}
