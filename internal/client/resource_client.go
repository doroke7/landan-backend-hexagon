package client

import (
	"google.golang.org/grpc"

	pbResourceModel "example/pb/resource/model"
)

func NewModel(oClientConn *grpc.ClientConn) *Model {
	return &Model{
		AdminUser: pbResourceModel.NewAdminUserClient(oClientConn),
	}
}

type Model struct {
	AdminUser pbResourceModel.AdminUserClient
}

////////////////////////////////////////////////////////////////////////////

func NewResourceClient(oClientConn *grpc.ClientConn, oModel *Model) *ResourceClient {

	return &ResourceClient{
		conn:  oClientConn,
		Model: oModel,
	}
}

type ResourceClient struct {
	conn   *grpc.ClientConn
	*Model // 這樣不厭其煩的命名 【嵌套結構】，是為了與 server 【命名空間一致性】，增加可讀性。
}

func (oClient *ResourceClient) Close() error {
	return oClient.conn.Close()
}
