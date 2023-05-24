package app

import (
	"bytes"
	"fmt"
	"log"
	"mime"

	"github.com/reviashko/remail/model"

	"github.com/knadh/go-pop3"
)

// POP3ClientInterface interface
type POP3ClientInterface interface {
	GetUnreadMessages(msgID int) ([]model.MessageInfo, error)
	Quit()
	PrintStat() error
}

// POP3Client struct
type POP3Client struct {
	Conn    *pop3.Conn
	Decoder *mime.WordDecoder
}

// NewPOP3Client func
func NewPOP3Client(server string, port int, tlsEnabled bool, login string, pswd string) POP3Client {

	p := pop3.New(pop3.Opt{
		Host:       server,
		Port:       port,
		TLSEnabled: tlsEnabled,
	})

	// Don't forget exec Quit before app close!
	c, err := p.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Auth(login, pswd); err != nil {
		log.Fatal(err)
	}

	dec := mime.WordDecoder{}

	return POP3Client{Conn: c, Decoder: &dec}
}

// Quit func
func (p *POP3Client) Quit() {
	p.Conn.Quit()
}

// PrintStat func
func (p *POP3Client) PrintStat() error {
	count, size, err := p.Conn.Stat()
	if err != nil {
		return err
	}
	fmt.Println("total messages=", count, "size=", size)

	return nil
}

// GetUnreadMessages func
func (p *POP3Client) GetUnreadMessages(msgID int) ([]model.MessageInfo, error) {

	retval := make([]model.MessageInfo, 0)
	dec := mime.WordDecoder{}

	msgs, _ := p.Conn.List(msgID)
	for _, m := range msgs {
		if m.ID <= msgID {
			continue
		}

		msg, _ := p.Conn.Retr(m.ID)

		subj, err := dec.DecodeHeader(msg.Header.Get("Subject"))
		if err != nil {
			//TODO: need to manage to this
			continue
		}

		from := msg.Header.Get("From")
		multiPart := msg.MultipartReader()

		// multi-part body:
		// https://github.com/emersion/go-message/blob/master/example_test.go#L12
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(msg.Body)
		if err != nil {
			log.Fatal(err)
		}

		retval = append(retval, model.MessageInfo{MsgID: m.ID, Subject: subj, IsMultiPart: multiPart != nil, From: from, Body: buf.Bytes()})
	}

	return retval, nil
}
