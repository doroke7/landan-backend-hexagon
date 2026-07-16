package interceptor_facade_admin

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"example/internal/helper"
)

type contextKey string

const AdminUserIDKey contextKey = "admin_user_id"

type AuthenticationInterceptor struct {
	*AbstractInterceptor
}

func NewAuthenticationInterceptor(oInterceptor *AbstractInterceptor) *AuthenticationInterceptor {
	return &AuthenticationInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *AuthenticationInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		oMd, bOk := metadata.FromIncomingContext(ctx)
		if !bOk {
			return nil, status.Error(codes.Unauthenticated, "缺少認證資訊")
		}

		aValues := oMd.Get("authorization")
		if len(aValues) == 0 {
			return nil, status.Error(codes.Unauthenticated, "缺少 authorization header")
		}

		sToken := strings.TrimPrefix(aValues[0], "Bearer ")

		oJwtHelper := helper.NewJwtHelper(helper.NewAbstractHelper())

		oClaims, oErr := oJwtHelper.Parse(sToken)
		if oErr != nil {
			return nil, status.Error(codes.Unauthenticated, "token 無效或已過期")
		}

		ctx = context.WithValue(ctx, AdminUserIDKey, oClaims.AdminUserId)
		return handler(ctx, req)
	}
}
