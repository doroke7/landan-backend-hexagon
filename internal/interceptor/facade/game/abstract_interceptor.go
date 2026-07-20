package interceptor_facade_admin

import (
	"context"
	helper "example/internal/helper"

	"google.golang.org/grpc"
)

type AbstractInterceptor struct {
	JwtHelper *helper.JwtHelper
}

func NewAbstractInterceptor() *AbstractInterceptor {
	return &AbstractInterceptor{
		JwtHelper: helper.NewJwtHelper(helper.NewAbstractHelper()),
	}
}

func (oSelf *AbstractInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(ctx, req)
	}
}
