package middleware_admin

import (
	"example/internal/bootstrap"
	"example/internal/utility"
	"net/url"
	"reflect"
	"unsafe"

	types "example/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"example/pkg"
)

type DecryptionMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewDecryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *DecryptionMiddleware {
	return &DecryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *DecryptionMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		// sHeaderKeys := oContext.GetHeader("Keys")
		sHeaderK := oContext.GetHeader("K")
		sHeaderA := oContext.GetHeader("A")

		sQueryS := oContext.DefaultQuery("s", "")
		sQueryO := oContext.DefaultQuery("o", "")
		sP := oContext.PostForm("p")

		var oRequestPayload types.RequestPayload

		// 2. 關鍵：使用 c.ShouldBind 代替 c.ShouldBindJSON！
		// Gin 會自動根據 Content-Type 去選用 JSON 解析器或 Form 解析器
		if err := oContext.ShouldBindBodyWith(&oRequestPayload, binding.JSON); err != nil {

			oContext.Abort()
			_ = oContext.Error(pkg.NewDefaultError("請求格式錯誤2", -1, 400))

			return
		}
		sP = oRequestPayload.P

		sKeys, oErr := oSelf.rsaHelper.Decrypt(sHeaderK, bootstrap.CONFIG.ADMIN.RSA.PRIVATE_KEY)
		if oErr != nil {
			oContext.Abort()
			_ = oContext.Error(pkg.NewDefaultError("金鑰解密失敗", -1, 400))
			return
		}

		// Go 的 encoding/json 只能反序列化到 exported（大写开头） 的字段。
		// 小写字段是 unexported 的，json.Unmarshal 无法访问，
		// 如果你写小写，直接跳过，解出来永远是空值。

		oKeys, _ := utility.JsonDecode[struct {
			Key string `json:"key"`
			Iv  string `json:"iv"`
		}](sKeys)

		oContext.Set("key", oKeys.Key)
		oContext.Set("iv", oKeys.Iv)

		sOption := oSelf.aesHelper.Decrypt(sQueryO, oKeys.Key, oKeys.Iv)

		sSearch := oSelf.aesHelper.Decrypt(sQueryS, oKeys.Key, oKeys.Iv)

		sParam := oSelf.aesHelper.Decrypt(sP, oKeys.Key, oKeys.Iv)

		sAuthorizaion := oSelf.aesHelper.Decrypt(sHeaderA, bootstrap.CONFIG.ADMIN.JWT.KEY, bootstrap.CONFIG.ADMIN.JWT.IV)
		oContext.Set("Authrization", sAuthorizaion)

		oOption, _ := utility.JsonDecode[struct {
			Size  string `json:"size"`
			Page  string `json:"page"`
			AppId string `json:"app_id"`
		}](sOption)

		oSearch, _ := utility.JsonDecode[map[string]interface{}](sSearch)
		oParam, _ := utility.JsonDecode[map[string]interface{}](sParam)

		oUrlQuery := oContext.Request.URL.Query()

		oSelf.Flatten(oUrlQuery, "search", oSearch)
		oSelf.Flatten(oUrlQuery, "option", oOption)
		oSelf.Flatten(oContext.Request.PostForm, "param", oParam)
		sEncoded := oUrlQuery.Encode()
		oContext.Request.URL.RawQuery = sEncoded

		/*
		   清理 Gin 的 queryCache / formCache，強制清理 context 中的參數緩存
		*/
		oContextValue := reflect.ValueOf(oContext).Elem()

		oQueryCacheField := oContextValue.FieldByName("queryCache")
		pQueryCacheField := oQueryCacheField.UnsafeAddr()
		ptrQueryCacheField := (*url.Values)(unsafe.Pointer(pQueryCacheField))
		*ptrQueryCacheField = nil

		oFormCacheField := oContextValue.FieldByName("formCache")
		pFormCacheField := oFormCacheField.UnsafeAddr()
		ptrFormCacheField := (*url.Values)(unsafe.Pointer(pFormCacheField))
		*ptrFormCacheField = nil

		oContext.Next()

	}
}
