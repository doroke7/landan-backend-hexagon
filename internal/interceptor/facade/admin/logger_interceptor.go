package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// AuthHolder 是塞進 context 往下傳的可變容器：LoggerInterceptor 在呼叫 handler 之前
// 建好、塞進 context，下游 handler（例如 SignIn）用 AuthHolderFromContext 把指標拿出來
// 直接寫進去；LoggerInterceptor 在 handler 呼叫「之後」讀 oHolder.Authorization 就讀得到——
// 因為兩邊指的是同一塊記憶體，資料是靠共享記憶體傳遞，不是靠 context.WithValue 本身往上傳
// （context 的值只能往下游流，這裡只是拿 context 夾帶指標）。
type AuthHolder struct {
	Authorization string
}

type authHolderKey struct{}

// AuthHolderFromContext 讓下游 handler 從 context 拿出 LoggerInterceptor 塞進去的 *AuthHolder。
func AuthHolderFromContext(oContext context.Context) (*AuthHolder, bool) {
	oHolder, bOk := oContext.Value(authHolderKey{}).(*AuthHolder)
	return oHolder, bOk
}

type LoggerInterceptor struct {
	*AbstractInterceptor
}

func NewLoggerInterceptor(oInterceptor *AbstractInterceptor) *LoggerInterceptor {
	return &LoggerInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *LoggerInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(
		oContext context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// before
		fmt.Println("Before LoggerInterceptor...")

		oHolder := &AuthHolder{}
		oContext = context.WithValue(oContext, authHolderKey{}, oHolder)

		oResponse, oReq := handler(oContext, req)
		fmt.Println("After LoggerInterceptor...")

		// after

		return oResponse, oReq
	}

}
