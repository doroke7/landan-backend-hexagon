package interceptor_facade_admin

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	bootstrap "example/bootstrap"
)

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

		pAuthrization := new(string)
		oContext = context.WithValue(oContext, "a", pAuthrization)

		oResponse, oErr = fnHandler(oContext, oRequest)
		if oErr != nil {
			return oResponse, oErr
		}

		if *pAuthrization == "" {
			return oResponse, oErr
		}

		sA := oSelf.AesHelper.Encrypt(
			*pAuthrization,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.KEY,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.IV,
		)

		if oHeaderErr := grpc.SetHeader(oContext, metadata.Pairs("A", sA)); oHeaderErr != nil {
			return nil, oHeaderErr
		}

		return oResponse, oErr
	}

}
