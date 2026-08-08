package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type EncryptionInterceptor struct {
	*AbstractInterceptor
}

func NewEncryptionInterceptor(oInterceptor *AbstractInterceptor) *EncryptionInterceptor {
	return &EncryptionInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 呼叫 handler 前先把 *authHolder 塞進 context，讓下游的 input 有地方寫入未加密
// 的 authorization；handler 執行完之後（after 階段）讀出這個明文，用 AES 加密，寫進
// "A" header——input 完全不碰加密邏輯，加密這件事統一由這個攔截器負責。
func (oSelf *EncryptionInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		sAuthorization := oContext.Value("authorization").(string)

		fmt.Println(sAuthorization)
		oResponse, oErr = fnHandler(oContext, oRequest)
		// if oErr != nil {
		// 	return oResponse, oErr
		// }

		// if sAuthorization != "" {
		// 	sA := oSelf.AesHelper.Encrypt(
		// 		sAuthorization,
		// 		bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.KEY,
		// 		bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.IV,
		// 	)

		// 	if oHeaderErr := grpc.SetHeader(oContext, metadata.Pairs("A", sA)); oHeaderErr != nil {
		// 		return nil, oHeaderErr
		// 	}
		// }

		return oResponse, oErr
	}

}
