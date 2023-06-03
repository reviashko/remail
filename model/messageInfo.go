package model

import (
	"bytes"
	"fmt"
)

// MessageInfo struct
type MessageInfo struct {
	MsgID   int
	Subject string
	From    string
	To      string
	Body    []byte
}

func (m *MessageInfo) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	buf.Write(m.Body)

	return buf.Bytes()
}
