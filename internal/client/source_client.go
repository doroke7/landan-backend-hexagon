package client

import (
	"google.golang.org/grpc"
	// pbResourceAnnouncement "example/pb/source/announcement"
)

func NewSourceClient(oClientConn *grpc.ClientConn) *SourceClient {

	return &SourceClient{
		conn: oClientConn,
	}
}

type SourceClient struct {
	conn *grpc.ClientConn
}

func (oClient *SourceClient) Close() error {
	return oClient.conn.Close()
}
