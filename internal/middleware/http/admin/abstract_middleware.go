package middleware_admin

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/gin-gonic/gin"

	pkg "example/pkg"

	helper "example/internal/helper"
)

type AbstractMiddleware struct {
	*pkg.Response
	rsaHelper *helper.RsaHelper
	aesHelper *helper.AesHelper
}

// 2. 在結構體上定義一個「構造函數」
func NewAbstractMiddleware(oResponse *pkg.Response, oRsaHelper *helper.RsaHelper, oAesHelper *helper.AesHelper) *AbstractMiddleware {
	return &AbstractMiddleware{
		Response:  oResponse,
		rsaHelper: oRsaHelper,
		aesHelper: oAesHelper,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *AbstractMiddleware) HandleAbstractMiddleware() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

	}
}

func (oSelf *AbstractMiddleware) Flatten(oQ url.Values, sPrefix string, mData interface{}) {
	val := reflect.ValueOf(mData)

	// 如果是指针，取其指向的实际内容
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		// 遍历 Map 的所有 Key
		for _, vKey := range val.MapKeys() {
			// 将 Key 转为字符串
			sKeyName := fmt.Sprintf("%v", vKey.Interface())
			sNewPrefix := sKeyName
			if sPrefix != "" {
				sNewPrefix = sPrefix + "." + sKeyName
			}
			// 递归处理子项
			oSelf.Flatten(oQ, sNewPrefix, val.MapIndex(vKey).Interface())
		}

	case reflect.Struct:
		// 遍历 Struct 的所有字段，使用 json tag 作为 key
		oType := val.Type()

		/*
			struct { Size string "json:\"size\""; Page string "json:\"page\""; AppId string "json:\"app_id\"" }
		*/
		// 把struct 的屬性for-each 處理
		for i := 0; i < val.NumField(); i++ {
			oField := oType.Field(i)
			sKeyName := oField.Tag.Get("json")
			if sKeyName == "" || sKeyName == "-" {
				sKeyName = oField.Name
			}
			sNewPrefix := sKeyName
			if sPrefix != "" {
				sNewPrefix = sPrefix + "." + sKeyName
			}

			oSelf.Flatten(oQ, sNewPrefix, val.Field(i).Interface())
		}

	case reflect.Slice, reflect.Array:
		// 处理数组，例如 search.items.0
		for i := 0; i < val.Len(); i++ {
			sIndexKey := fmt.Sprintf("%s.%d", sPrefix, i)
			oSelf.Flatten(oQ, sIndexKey, val.Index(i).Interface())
		}

	default:
		// 递归终点：将基础类型转为 string 存入 url.Values
		// 即使 mData 原本是 int，存入 Query 时也必须是 string
		oQ.Set(sPrefix, fmt.Sprintf("%v", val.Interface()))
	}
}
