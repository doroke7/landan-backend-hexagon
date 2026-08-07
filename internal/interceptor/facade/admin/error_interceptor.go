package interceptor_facade_admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorInterceptor struct {
	*AbstractInterceptor
}

func NewErrorInterceptor(oInterceptor *AbstractInterceptor) *ErrorInterceptor {
	return &ErrorInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *ErrorInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(oContext context.Context, oReqeust any, oServerIno *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		oErr = nil

		fmt.Println("Before ErrorInterceptor...")

		defer func() {

			if oPanic := recover(); oPanic != nil {
				oResponse = nil
				oErr = status.Error(codes.Unavailable, "系統錯誤")
			}
		}()

		oResponse, oErr = fnHandler(oContext, oReqeust)
		fmt.Println("After ErrorInterceptor...")

		// after

		return oResponse, oErr
	}

}
