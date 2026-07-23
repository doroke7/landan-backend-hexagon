package types

type RequestPayload struct {
	P string `json:"p" form:"p" binding:"required"`
}

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
