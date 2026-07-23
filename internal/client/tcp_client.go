package client

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"

	types "example/types"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

func NewTcpClient(oConn net.Conn) *TcpClient {
	return &TcpClient{
		Conn:   oConn,
		reader: bufio.NewReader(oConn),
	}
}

// reader 在整條連線的生命週期內只建一次，所有方法共用同一個 *bufio.Reader，
// 避免每次呼叫都重建：重建的話，前一次可能已經多讀進緩衝區、還沒吐出來的
// 下一筆 response 的字節會跟著舊的 reader 一起被丟掉，造成連線從此對不上帧。
type TcpClient struct {
	net.Conn
	reader *bufio.Reader
}

func (oSelf *TcpClient) AdminAuthenticationAuthenticatorSignIn(sName string, sPassword string) (*types.TcpResponse, error) {

	aFrame, err := oSelf.EncodeFrame(types.TcpRequest{
		Method: "Admin.Authentication.Authenticator.SignIn",
		Param:  sName + ":" + sPassword,
	})
	if err != nil {
		log.Println("tcp client: encode failed:", err)
		return nil, err
	}

	if _, err = oSelf.Conn.Write(aFrame); err != nil {
		log.Println("tcp client: write failed:", err)
		return nil, err
	}

	var oResp types.TcpResponse
	// IMPORTANT
	// TCP 是stream， 你有可能讀取 到 第二個 response 的資料， oSelf.DecodeFrame(bufio.NewReader(oSelf.Conn), &oResp)
	if err = oSelf.DecodeFrame(oSelf.reader, &oResp); err != nil {
		log.Println("tcp client: decode failed:", err)
		return nil, err
	}

	return &oResp, nil
}

func (oSelf *TcpClient) EncodeFrame(oPayload any) ([]byte, error) {
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

func (oSelf *TcpClient) DecodeFrame(oReader *bufio.Reader, oPayload any) error {
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
