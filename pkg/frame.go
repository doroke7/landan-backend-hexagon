package pkg

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

/*
TCP 是位元組流（stream），沒有天然的訊息邊界，直接 conn.Read(buf) 會遇到：
  - 黏包：對方連續送兩筆小訊息，一次 Read 可能把兩筆都讀進來
  - 拆包：一筆訊息太大，一次 Read 讀不完，要跨好幾次 Read 才拼得齊

解法是自訂一個「帶長度前綴」的 frame 格式，讀取端固定「先讀長度、再讀滿長度」；
中間跨了幾次底層 Read 都由 bufio.Reader + io.ReadFull 自動處理，
呼叫端每呼叫一次 ReadFrame 就保證拿到剛好一個完整的 method + message。

Frame 格式（Big Endian）：

	4 bytes   BodyLength    body 的長度（不含這 4 bytes 自己）
	1 byte    MethodLength  method 名稱的長度
	N bytes   Method        方法名稱，例如 "Admin.Authentication.SignIn"
	M bytes   Message       實際訊息內容，M = BodyLength - 1 - N

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
)

// EncodeFrame 把 method + message 封裝成一個完整的 frame，寫出去的一方（server 回包／client 發送）都用這個。
func EncodeFrame(sMethod string, aMessage []byte) ([]byte, error) {
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

/*

整体思路：TCP 是字节流,没有天然的消息边界,所以自定义一个「长度前缀」frame 格式,读的一方先知道这一包多长,再照着长度读满,天然解决黏包/拆包问题。

格式(Big Endian)：

┌──────────────┬──────────────┬─────────────┬──────────────┐
│  BodyLength  │ MethodLength │   Method    │   Message    │
│   4 bytes    │    1 byte    │   N bytes   │   M bytes    │
└──────────────┴──────────────┴─────────────┴──────────────┘
                └──────────────── BodyLength ─────────────┘


*/

// ReadFrame 從 reader 讀出「剛好一個」完整 frame，讀不滿會阻塞等下一批 bytes，
// 不會因為單次底層 Read 只讀到半包、或一次讀到好幾包而出錯——這就是拆包/黏包的解法。
func ReadFrame(oReader *bufio.Reader) (sMethod string, aMessage []byte, err error) {

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
