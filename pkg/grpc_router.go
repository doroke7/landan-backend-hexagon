package pkg

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

/*
GrpcRouter 模仿 Hyperf 的 Router::addGroup(prefix, ['middleware' => [...]])：
group 註冊時就把攔截器綁死在 prefix 上，dispatch 時只需要對 gRPC 的
info.FullMethod（格式為 "/package.Service/Method"）做前綴比對，
不再靠切字串位置去猜 group。

aBase 是全局攔截器（等同 Hyperf 的全局 middleware），一定會套用；
Group 註冊的攔截器會疊加在 aBase 之上。沒有任何 Group 匹配時，
仍然套用 aBase，不會讓請求完全繞過 error/log 等防護。
*/

type Route struct {
	Prefix       string
	Interceptors []grpc.UnaryServerInterceptor
}

type GrpcRouter struct {
	base   []grpc.UnaryServerInterceptor
	routes []Route
}

func NewGrpcRouter(aBase ...grpc.UnaryServerInterceptor) *GrpcRouter {
	return &GrpcRouter{base: aBase}
}

func ChainInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			i, prev := i, chain
			chain = func(c context.Context, r any) (any, error) {
				return interceptors[i](c, r, info, prev)
			}
		}
		return chain(ctx, req)
	}
}

// Group 註冊一個 prefix 對應的攔截器鏈，aExtra 會疊加在 aBase 之後
func (oSelf *GrpcRouter) Group(sPrefix string, aExtra ...grpc.UnaryServerInterceptor) *GrpcRouter {
	aInterceptors := make([]grpc.UnaryServerInterceptor, 0, len(oSelf.base)+len(aExtra))
	aInterceptors = append(aInterceptors, oSelf.base...)
	aInterceptors = append(aInterceptors, aExtra...)

	oSelf.routes = append(oSelf.routes, Route{
		Prefix:       sPrefix,
		Interceptors: aInterceptors,
	})
	return oSelf
}

func (oSelf *GrpcRouter) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		for _, oRoute := range oSelf.routes {
			if strings.HasPrefix(info.FullMethod, oRoute.Prefix) {
				return ChainInterceptors(oRoute.Interceptors...)(ctx, req, info, handler)
			}
		}
		// 沒有任何 group 匹配到，至少套用全局攔截器，而不是完全放行
		return ChainInterceptors(oSelf.base...)(ctx, req, info, handler)
	}
}
