package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// MessageInfo struct
type MessageInfo struct {
	MsgID   int
	Subject string
	From    string
	To      string
	Body    []byte
}

type MessageMP struct {
	ContentTye string
	Headers    map[string]string
	Body       []byte
}

type Message struct {
	MsgID           int
	BodyContentType string
	From            string
	To              []string
	CC              []string
	BCC             []string
	Subject         string
	Body            []byte
	MultiParts      []MessageMP
	Attachments     map[string][]byte
}

func NewMessage(subject string) *Message {
	return &Message{Subject: subject, MultiParts: make([]MessageMP, 0), Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) AttachBase64(fileName string, b []byte) error {

	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	multiParts := len(m.MultiParts) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments || multiParts {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	if multiParts {
		log.Println("Multiparts")

		for _, v := range m.MultiParts {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))

			for k, h := range v.Headers {
				buf.WriteString(fmt.Sprintf("%s: %s\n", k, h))

				log.Println(fmt.Sprintf("%s: %s\n", k, h))
			}

			if strings.Contains(v.ContentTye, "octet-stream") {
				//b := make([]byte, base64.StdEncoding.EncodedLen(len(v.Body)))
				//base64.StdEncoding.Encode(b, v.Body)
				buf.Write(v.Body)
			} else {
				buf.Write(v.Body)
			}

			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
	} else {
		buf.Write(m.Body)
	}

	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
