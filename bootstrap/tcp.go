package bootstrap

import (
	"fmt"
	"net"
	"time"
)

// NewTcp 連到自家 tcp 服務（cmd/tcp.go 開的那個），角色跟 NewResource 一樣，
// 從 CONFIG.CLIENTS.TCP.HOSTS/PORTS 建立連線。
// 跟 gRPC 的 resource client 不同，這裡沒有 resolver／round-robin 那套機制，
// 一個 net.Conn 就是一條實體 TCP 連線；多節點時依序嘗試，連上第一個就回傳。
func NewTcp() (net.Conn, error) {
	aHosts := CONFIG.CLIENTS.TCP.HOSTS
	aPorts := CONFIG.CLIENTS.TCP.PORTS

	oDialer := net.Dialer{
		Timeout:   3 * time.Second,  // 建立連線最多等 3 秒
		KeepAlive: 10 * time.Second, // 底層 socket 保活，每 10 秒探測一次，跟 resource client 的 keepalive 同樣目的
	}

	var err error
	for iIndex, sHost := range aHosts {
		sPort := aPorts[0]
		if iIndex < len(aPorts) {
			sPort = aPorts[iIndex]
		}

		var oConn net.Conn
		if oConn, err = oDialer.Dial("tcp", fmt.Sprintf("%s:%s", sHost, sPort)); err == nil {
			return oConn, nil
		}
	}

	return nil, err
}
