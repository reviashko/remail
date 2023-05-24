package model

// MessageInfo struct
type MessageInfo struct {
	MsgID       int
	Subject     string
	From        string
	IsMultiPart bool
	Body        []byte
}
