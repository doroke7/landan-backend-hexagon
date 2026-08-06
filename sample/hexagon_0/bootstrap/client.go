package bootstrap

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(dsn string) (*grpc.ClientConn, error) {
	return grpc.NewClient(dsn, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
