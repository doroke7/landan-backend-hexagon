package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type EncrtptionInterceptor struct {
	*AbstractInterceptor
}

func NewEncrtptionInterceptor(oInterceptor *AbstractInterceptor) *EncrtptionInterceptor {
	return &EncrtptionInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *EncrtptionInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(oContext context.Context, oReqeust any, oServerIno *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		oErr = nil

		fmt.Println("Before EncrtptionInterceptor...")
		oResponse, oErr = fnHandler(oContext, oReqeust)
		fmt.Println("After EncrtptionInterceptor...")

		oMd, _ := metadata.FromIncomingContext(oContext)
		_ = oMd // TODO: 還沒接上回應加密邏輯

		return oResponse, oErr
	}

}
