package app

import (
	"bytes"

	"github.com/knadh/go-pop3"
	"github.com/reviashko/remail/model"
)

// POP3ClientInterface interface
type POP3ClientInterface interface {
	GetUnreadMessages(msgID int) ([]model.MessageInfo, error)
}

// POP3Client struct
type POP3Client struct {
	Client *pop3.Client
	Login  string
	Pswd   string
}

// NewPOP3Client func
func NewPOP3Client(server string, port int, tlsEnabled bool, login string, pswd string) POP3Client {

	p := pop3.New(pop3.Opt{
		Host:       server,
		Port:       port,
		TLSEnabled: tlsEnabled,
	})

	return POP3Client{Client: p, Login: login, Pswd: pswd}
}

// GetUnreadMessages func
func (p *POP3Client) GetUnreadMessages(msgID int) ([]model.MessageInfo, error) {

	retval := make([]model.MessageInfo, 0)

	connect, err := p.Client.NewConn()
	if err != nil {
		return retval, err
	}
	defer connect.Quit()

	if err := connect.Auth(p.Login, p.Pswd); err != nil {
		return retval, err
	}

	// TODO: need to avoid reading all messages
	msgs, _ := connect.List(0)
	for _, m := range msgs {
		if m.ID <= msgID {
			continue
		}

		msg, _ := connect.Retr(m.ID)

		from := msg.Header.Get("From")
		to := msg.Header.Get("To")

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(msg.Body)
		if err != nil {
			return retval, err
		}

		retval = append(retval, model.MessageInfo{MsgID: m.ID, To: to, From: from, Body: buf.Bytes()})
	}

	return retval, nil
}
