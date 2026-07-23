package pkg

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
)

// TcpHandlerFunc 是一個 method 對應的處理方法，簽名統一，方便用 method name 當 key 做路由。
// 拆包/黏包已經在 ReadFrame／EncodeFrame 處理過，這裡拿到的都是完整的一個 message，
// 回傳的 []byte 會被包成 frame 寫回連線。
type TcpHandlerFunc func(aMessage []byte) ([]byte, error)

var ErrTcpMethodNotFound = errors.New("tcp: method not found")

// TcpRouter 職責跟 ConsumerRouter 一樣：只負責把 method name 對應到一個處理方法，
// 不管 unmarshal／business 邏輯；Serve 時自己負責 accept 連線、讀 frame、分發、回包。
type TcpRouter struct {
	routes map[string]TcpHandlerFunc
}

func NewTcpRouter() *TcpRouter {
	return &TcpRouter{routes: make(map[string]TcpHandlerFunc)}
}

// HandleFunc 註冊一個 method 對應的處理方法，用法跟 ConsumerRouter.HandleFunc 一樣。
func (oSelf *TcpRouter) HandleFunc(sMethod string, fnHandler TcpHandlerFunc) *TcpRouter {
	oSelf.routes[sMethod] = fnHandler
	return oSelf
}

// Serve 監聽 sAddr，每個連線各自在自己的 goroutine 讀 frame、分發、回包，
// ctx 取消時關掉 listener，讓 Accept 中斷、Serve 返回。
func (oSelf *TcpRouter) Serve(ctx context.Context, sAddr string) error {
	oListener, err := net.Listen("tcp", sAddr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		oListener.Close()
	}()

	for {
		oConn, err := oListener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}

		go oSelf.serveConn(oConn)
	}
}

func (oSelf *TcpRouter) serveConn(oConn net.Conn) {
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	for {
		sMethod, aMessage, err := ReadFrame(oReader)
		if err != nil {
			return
		}

		aResp, err := oSelf.dispatch(sMethod, aMessage)
		if err != nil {
			log.Printf("tcp: dispatch failed: method=%s err=%v", sMethod, err)
			continue
		}

		aFrame, err := EncodeFrame(sMethod, aResp)
		if err != nil {
			log.Printf("tcp: encode failed: method=%s err=%v", sMethod, err)
			continue
		}

		if _, err = oConn.Write(aFrame); err != nil {
			return
		}
	}
}

func (oSelf *TcpRouter) dispatch(sMethod string, aMessage []byte) ([]byte, error) {
	fnHandler, ok := oSelf.routes[sMethod]
	if !ok {
		return nil, ErrTcpMethodNotFound
	}
	return fnHandler(aMessage)
}
