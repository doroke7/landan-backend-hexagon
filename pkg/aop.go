package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Aop 結構體：不帶泛型，是一個純粹、萬用的切面工具箱
type Aop struct {
	oCache redis.UniversalClient
}

// NewAop 初始化時注入萬用的 Redis 客戶端
func NewAop(oCache redis.UniversalClient) *Aop {
	return &Aop{oCache: oCache}
}

// Cacheable 順利綁定為 Aop 結構體的方法！
// 💡 關鍵設計：
// 1. pDest: 傳入用來接收結果的「結構體指針」（如同 json.Unmarshal 的第二個參數）
// 2. cFn: 核心業務函數（如查 DB），回傳萬用的 interface{}
func (oSelf *Aop) Cacheable(oCtx context.Context, sKey string, oTtl time.Duration, pDest interface{}, cFn func(oCtx context.Context) (interface{}, error)) error {

	// 1. 前置安全檢查
	if oErr := oCtx.Err(); oErr != nil {
		return oErr
	}

	// 2. 前置切面 (Before)：查快取
	sJsonStr, oErr := oSelf.oCache.Get(oCtx, sKey).Result()
	if oErr == nil {
		// 🎉 快取命中 (Hit)！直接反序列化進傳進來的 pDest 指針
		if oErr := json.Unmarshal([]byte(sJsonStr), pDest); oErr == nil {
			fmt.Printf("🎯 [Aop] 完美命中快取！資料已自動注入。Key: '%s'\n", sKey)
			return nil
		}
	} else if !errors.Is(oErr, redis.Nil) {
		// 快取降級：Redis 異常不中斷主業務，讓流量穿透去查 DB
		fmt.Printf("⚠️ [Aop] 快取連線異常: %v，自動降級穿透。\n", oErr)
	}

	fmt.Printf("🔍 [Aop] 快取未命中 (Miss)，準備執行核心業務... Key: '%s'\n", sKey)

	// 3. 執行核心業務 (Join Point：例如查 MySQL)
	oRes, oErr := cFn(oCtx)
	if oErr != nil {
		return oErr
	}

	// 4. 後置切面 (After Returning)：將真實業務回傳的物件，序列化同步到 Redis
	aByteData, oErr := json.Marshal(oRes)
	if oErr == nil {
		oSelf.oCache.Set(oCtx, sKey, string(aByteData), oTtl)
		fmt.Printf("💾 [Aop] 數據已自動同步至 Redis（TTL: %v）。Key: '%s'\n", oTtl, sKey)
	}

	// 5. 💡 核心魔法：如果快取沒中，查完 DB 後，要把結果深拷貝（Deep Copy）給外面的 pDest 指針
	if oErr := json.Unmarshal(aByteData, pDest); oErr != nil {
		return oErr
	}

	return nil
}
