package bootstrap

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func NewClient() (*grpc.ClientConn, error) {
	aHosts := CONFIG.CLIENTS.FACADE.HOSTS
	aPorts := CONFIG.CLIENTS.FACADE.PORTS

	aAddresses := make([]resolver.Address, 0, len(aHosts))
	for i, sHost := range aHosts {
		var sPort string
		if i < len(aPorts) {
			sPort = aPorts[i]
		}
		aAddresses = append(aAddresses, resolver.Address{Addr: fmt.Sprintf("%s:%s", sHost, sPort)})
	}

	// manual resolver 把 CLIENTS.FACADE.HOSTS/PORTS 這組節點清單，一次性交給同一個
	// *grpc.ClientConn；搭配 round_robin，client 會同時跟每個節點保持連線、輪詢送出
	// 請求，某個節點斷線時 grpc 會自動跳過它，不用自己手動挑 host 或做 failover。
	oBuilder := manual.NewBuilderWithScheme("facade")
	oBuilder.InitialState(resolver.State{Addresses: aAddresses})

	return grpc.NewClient(
		oBuilder.Scheme()+":///",
		grpc.WithResolvers(oBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
}
