package interceptor_facade_admin

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	bootstrap "example/bootstrap"
)

// AuthHolderFromContext 讓下游的 SignIn 從 context 拿出 EncryptionInterceptor 塞進去的
// 指標，直接寫入未加密的 authorization——EncryptionInterceptor 呼叫 handler 前建好、
// 塞進 context，handler 執行完之後讀 oHolder.Authorization 就讀得到，兩邊讀寫的是同一
// 塊記憶體，不是靠 context 本身把值傳回來。用匿名 struct，不額外宣告具名型別。
func AuthHolderFromContext(oContext context.Context) (*struct{ Authorization string }, bool) {
	oHolder, bOk := oContext.Value("a").(*struct{ Authorization string })
	return oHolder, bOk
}

type EncryptionInterceptor struct {
	*AbstractInterceptor
}

func NewEncryptionInterceptor(oInterceptor *AbstractInterceptor) *EncryptionInterceptor {
	return &EncryptionInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 呼叫 handler 前把 *authHolder 塞進 context；handler 執行完之後（after 階段）
// 讀出 SignIn 寫進 holder 的明文 authorization，用 AES 加密，寫進 "A" header。
func (oSelf *EncryptionInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		oHolder := &struct{ Authorization string }{}
		oContext = context.WithValue(oContext, "a", oHolder)

		oResponse, oErr = fnHandler(oContext, oRequest)
		if oErr != nil {
			return oResponse, oErr
		}

		if oHolder.Authorization == "" {
			return oResponse, oErr
		}

		sA := oSelf.AesHelper.Encrypt(
			oHolder.Authorization,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.KEY,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.IV,
		)

		if oHeaderErr := grpc.SetHeader(oContext, metadata.Pairs("A", sA)); oHeaderErr != nil {
			return nil, oHeaderErr
		}

		return oResponse, oErr
	}

}
