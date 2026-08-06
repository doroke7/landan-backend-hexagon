package client

import (
	"google.golang.org/grpc"

	pbResourceAnnouncement "example/pb/source/announcement"
)

func NewAnnouncement(oClientConn *grpc.ClientConn) *Announcement {

	return &Announcement{
		Lottery: pbResourceAnnouncement.NewLotteryClient(oClientConn),
	}
}

type Announcement struct {
	Lottery pbResourceAnnouncement.LotteryClient
}

func NewSourceClient(oClientConn *grpc.ClientConn, oAnnouncement *Announcement) *SourceClient {

	return &SourceClient{
		conn:         oClientConn,
		Announcement: oAnnouncement,
	}
}

type SourceClient struct {
	conn         *grpc.ClientConn
	Announcement *Announcement
}

func (oClient *SourceClient) Close() error {
	return oClient.conn.Close()
}
