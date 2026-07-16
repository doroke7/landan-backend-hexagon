package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type StatusInterceptor struct {
	*AbstractInterceptor
}

func NewStatusInterceptor(oInterceptor *AbstractInterceptor) *StatusInterceptor {
	return &StatusInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *StatusInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// before
		fmt.Println("Before StatusInterceptor...")

		oCtx, oReq := handler(ctx, req)
		fmt.Println("After StatusInterceptor...")

		// after

		return oCtx, oReq
	}

}
