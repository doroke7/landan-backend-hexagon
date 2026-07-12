package bootstrap

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient() (*grpc.ClientConn, error) {
	var sHost, sPort string
	if len(CONFIG.CLIENTS.FACADE.HOSTS) > 0 {
		sHost = CONFIG.CLIENTS.FACADE.HOSTS[0]
	}
	if len(CONFIG.CLIENTS.FACADE.PORTS) > 0 {
		sPort = CONFIG.CLIENTS.FACADE.PORTS[0]
	}
	sAddr := fmt.Sprintf("%s:%s", sHost, sPort)

	return grpc.NewClient(sAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
