package client

import (
	"google.golang.org/grpc"

	pbClient "example/pb/client"
)

type Client struct {
	User pbClient.UserServiceClient
}

func NewClient(oConn *grpc.ClientConn) *Client {
	return &Client{
		User: pbClient.NewUserServiceClient(oConn),
	}
}
