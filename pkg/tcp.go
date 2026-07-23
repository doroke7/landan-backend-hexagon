package pkg

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

/*
TCP 是位元組流（stream），沒有天然的訊息邊界，直接 conn.Read(buf) 會遇到：
  - 黏包：對方連續送兩筆小訊息，一次 Read 可能把兩筆都讀進來
  - 拆包：一筆訊息太大，一次 Read 讀不完，要跨好幾次 Read 才拼得齊

解法是自訂一個「帶長度前綴」的 frame 格式，讀取端固定「先讀長度、再讀滿長度」；
中間跨了幾次底層 Read 都由 bufio.Reader + io.ReadFull 自動處理，
呼叫端每呼叫一次 ReadFrame 就保證拿到剛好一個完整的 method + message。

Frame 格式（Big Endian）：

	┌──────────────┬──────────────┬─────────────┬──────────────┐
	│  BodyLength  │ MethodLength │   Method    │   Message    │
	│   4 bytes    │    1 byte    │   N bytes   │   M bytes    │
	└──────────────┴──────────────┴─────────────┴──────────────┘
	                └──────────────── BodyLength ─────────────┘

method 用顯式長度前綴切割，而不是用分隔符（例如空白或 \n）去切，
是因為 Message 內容可能是任意 binary，用分隔符切會有跟 payload 內容衝突的風險。
*/

const (
	tcpMaxBodyLength   = 1 << 12  // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
	tcpMaxMethodLength = 1<<8 - 1 // MethodLength 只有 1 byte，最大 255
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodTooLarge = errors.New("tcp: method name too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

// 回傳的 []byte 會被包成 frame 寫回連線。
type TcpHandlerFunc func(aMessage []byte) ([]byte, error)

// Tcp 職責跟 ConsumerRouter 一樣：只負責把 method name 對應到一個處理方法，
// 不管 unmarshal／business 邏輯；Serve 時自己負責 accept 連線、讀 frame、分發、回包。
type Tcp struct {
	routes map[string]TcpHandlerFunc
}

func NewTcp() *Tcp {
	return &Tcp{routes: make(map[string]TcpHandlerFunc)}
}

// HandleFunc 註冊一個 method 對應的處理方法，用法跟 ConsumerRouter.HandleFunc 一樣。
func (oSelf *Tcp) HandleFunc(sMethod string, fnHandler TcpHandlerFunc) *Tcp {
	oSelf.routes[sMethod] = fnHandler
	return oSelf
}

// Serve 監聽 sAddr，每個連線各自在自己的 goroutine 讀 frame、分發、回包，
// ctx 取消時關掉 listener，讓 Accept 中斷、Serve 返回。
func (oSelf *Tcp) Serve(ctx context.Context, sAddr string) error {
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

func (oSelf *Tcp) serveConn(oConn net.Conn) {
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	for {
		sMethod, aMessage, err := oSelf.DecodeFrame(oReader)
		if err != nil {
			return
		}

		aResp, err := oSelf.dispatch(sMethod, aMessage)
		if err != nil {
			log.Printf("tcp: dispatch failed: method=%s err=%v", sMethod, err)
			continue
		}

		aFrame, err := oSelf.EncodeFrame(sMethod, aResp)
		if err != nil {
			log.Printf("tcp: encode failed: method=%s err=%v", sMethod, err)
			continue
		}

		if _, err = oConn.Write(aFrame); err != nil {
			return
		}
	}
}

func (oSelf *Tcp) dispatch(sMethod string, aMessage []byte) ([]byte, error) {
	fnHandler, ok := oSelf.routes[sMethod]
	if !ok {
		return nil, ErrTcpMethodNotFound
	}
	return fnHandler(aMessage)
}

// EncodeFrame 把 method + message 封裝成一個完整的 frame，寫出去的一方（server 回包／client 發送）都用這個。
func (oSelf *Tcp) EncodeFrame(sMethod string, aMessage []byte) ([]byte, error) {
	if len(sMethod) > tcpMaxMethodLength {
		return nil, ErrTcpMethodTooLarge
	}

	iBodyLength := 1 + len(sMethod) + len(aMessage)
	if iBodyLength > tcpMaxBodyLength {
		return nil, ErrTcpBodyTooLarge
	}

	aFrame := make([]byte, 4+iBodyLength)
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(iBodyLength))
	aFrame[4] = byte(len(sMethod))
	copy(aFrame[5:5+len(sMethod)], sMethod)
	copy(aFrame[5+len(sMethod):], aMessage)

	return aFrame, nil
}

func (oSelf *Tcp) DecodeFrame(oReader *bufio.Reader) (sMethod string, aMessage []byte, err error) {

	aLengthBuf := make([]byte, 4)
	if _, err = io.ReadFull(oReader, aLengthBuf); err != nil {
		return "", nil, err
	}

	iBodyLength := binary.BigEndian.Uint32(aLengthBuf)
	if iBodyLength == 0 || iBodyLength > tcpMaxBodyLength {
		return "", nil, ErrTcpBodyTooLarge
	}

	aBody := make([]byte, iBodyLength)
	if _, err = io.ReadFull(oReader, aBody); err != nil {
		return "", nil, err
	}

	iMethodLength := int(aBody[0])
	if len(aBody) < 1+iMethodLength {
		return "", nil, ErrTcpMethodTooLarge
	}

	sMethod = string(aBody[1 : 1+iMethodLength])
	aMessage = aBody[1+iMethodLength:]

	return sMethod, aMessage, nil
}
