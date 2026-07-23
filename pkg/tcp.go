package pkg

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
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
中間跨了幾次底層 Read 都由 bufio.Reader + io.ReadFull 自動處理。

Frame 格式（Big Endian）：

	┌──────────────┬─────────────────────────┐
	│  BodyLength  │   Body (JSON)           │
	│   4 bytes    │   N bytes               │
	└──────────────┴─────────────────────────┘

Body 是 JSON，內容依方向不同：
  - Request  ：{ code, method, param }
  - Response ：{ code, message, result }
*/

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

// TcpRequest 是 client 送給 server 的內容，method 用來給 Tcp 分發到對應的 handler。
type TcpRequest struct {
	Code   int    `json:"code"`
	Method string `json:"method"`
	Param  string `json:"param"`
}

// TcpResponse 是 server 回給 client 的內容，跟 pkg.Response 的 code/message/result 是同一套慣例。
type TcpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

// TcpHandlerFunc 是一個 method 對應的處理方法，簽名統一，方便用 method name 當 key 做路由。
type TcpHandlerFunc func(oReq TcpRequest) TcpResponse

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
		oReq, err := oSelf.DecodeRequestFrame(oReader)
		if err != nil {
			return
		}

		oResp := oSelf.dispatch(oReq)

		aFrame, err := oSelf.EncodeResponseFrame(oResp)
		if err != nil {
			log.Printf("tcp: encode failed: method=%s err=%v", oReq.Method, err)
			continue
		}

		if _, err = oConn.Write(aFrame); err != nil {
			return
		}
	}
}

func (oSelf *Tcp) dispatch(oReq TcpRequest) TcpResponse {
	fnHandler, ok := oSelf.routes[oReq.Method]
	if !ok {
		return TcpResponse{Code: -1, Message: ErrTcpMethodNotFound.Error()}
	}
	return fnHandler(oReq)
}

// EncodeRequestFrame 是 client 端組 request 用的。
func (oSelf *Tcp) EncodeRequestFrame(oReq TcpRequest) ([]byte, error) {
	return oSelf.encodeFrame(oReq)
}

// EncodeResponseFrame 是 server 端組 response 用的。
func (oSelf *Tcp) EncodeResponseFrame(oResp TcpResponse) ([]byte, error) {
	return oSelf.encodeFrame(oResp)
}

// DecodeRequestFrame 是 server 端讀 client request 用的。
func (oSelf *Tcp) DecodeRequestFrame(oReader *bufio.Reader) (TcpRequest, error) {
	var oReq TcpRequest
	err := oSelf.decodeFrame(oReader, &oReq)
	return oReq, err
}

// DecodeResponseFrame 是 client 端讀 server response 用的。
func (oSelf *Tcp) DecodeResponseFrame(oReader *bufio.Reader) (TcpResponse, error) {
	var oResp TcpResponse
	err := oSelf.decodeFrame(oReader, &oResp)
	return oResp, err
}

func (oSelf *Tcp) encodeFrame(oPayload any) ([]byte, error) {
	aBody, err := json.Marshal(oPayload)
	if err != nil {
		return nil, err
	}

	if len(aBody) > tcpMaxBodyLength {
		return nil, ErrTcpBodyTooLarge
	}

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	return aFrame, nil
}

func (oSelf *Tcp) decodeFrame(oReader *bufio.Reader, oPayload any) error {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return err
	}

	iBodyLength := binary.BigEndian.Uint32(aLengthBuf)
	if iBodyLength == 0 || iBodyLength > tcpMaxBodyLength {
		return ErrTcpBodyTooLarge
	}

	aBody := make([]byte, iBodyLength)
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return err
	}

	return json.Unmarshal(aBody, oPayload)
}
