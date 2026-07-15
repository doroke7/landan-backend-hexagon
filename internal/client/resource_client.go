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

/*
    為何這裡寫 【AdminUser pbResourceModel.AdminUserClient】 而不是 【AdminUser *pbResourceModel.AdminUserClient】
	是可行的？
	1. 首先寫 * 的用意是為了 全系統 連線變量唯一
	2. 何時可能不唯一 -> 存在兩個地方使用， 譬如 Controller 跟 Middleware 都注入了 Helper類
	3. 但是這邊 AdminUserClient 只有 Resource.Model 使用


	4. AdminUserClient 是 Interface， 我們不需要對 interface 寫 *
*/

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
