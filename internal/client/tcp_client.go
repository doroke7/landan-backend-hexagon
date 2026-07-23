package client

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
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
		Conn: oConn,
	}
}

type TcpClient struct {
	net.Conn
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
