package types

type RequestPayload struct {
	P string `json:"p" form:"p" binding:"required"`
}
