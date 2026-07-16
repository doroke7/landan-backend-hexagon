package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

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
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// before
		fmt.Println("Before LoggerInterceptor...")

		oCtx, oReq := handler(ctx, req)
		fmt.Println("After LoggerInterceptor...")

		// after

		return oCtx, oReq
	}

}
