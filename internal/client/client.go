package client

import (
	"google.golang.org/grpc"

	pb "example/pb/client"
)

type Client struct {
	User pb.UserServiceClient
}

func NewClient(oConn *grpc.ClientConn) *Client {
	return &Client{
		User: pb.NewUserServiceClient(oConn),
	}
}
