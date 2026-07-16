package interceptor_resource

import (
	"context"
	"strings"

	bootstrap "example/internal/bootstrap"
	utility "example/internal/utility"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AllInterceptor struct {
	*AbstractInterceptor
}

func NewAllInterceptor(oInterceptor *AbstractInterceptor) *AllInterceptor {
	return &AllInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 驗證 facade -> resource 呼叫的 Basic Auth
func (oSelf *AllInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		oMetadata, bMetaOb := metadata.FromIncomingContext(ctx)

		if bMetaOb {
			aAuthotizations := oMetadata.Get("authorization")
			sAuthotizations := strings.Join(aAuthotizations, "")
			sName := bootstrap.CONFIG.SERVICES.RESOURCE.NAME
			sPassword := bootstrap.CONFIG.SERVICES.RESOURCE.PASSWORD

			sAuthotization := "Basic " + utility.Base64Encode(sName+":"+sPassword)

			if sAuthotizations != sAuthotization {
				return nil, status.Error(codes.PermissionDenied, "帳號密碼錯誤")
			}
		}

		return handler(ctx, req)
	}
}
