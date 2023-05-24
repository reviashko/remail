package app

import (
	"fmt"
	"net/smtp"
)

// SMTPClientInterface interface
type SMTPClientInterface interface {
	SendEmail(toEmails []string, message []byte) error
}

// SMTPClient struct
type SMTPClient struct {
	Auth   *smtp.Auth
	Port   int
	Server string
	Login  string
}

// NewSMTPClient func
// login := "from@email.host"
// pswd := "email password"
// server := "smtp.email.host"
// port := 587
func NewSMTPClient(server string, port int, login string, pswd string) SMTPClient {

	auth := smtp.PlainAuth("", login, pswd, server)
	return SMTPClient{Auth: &auth, Port: port, Server: server, Login: login}
}

// SendEmail func
func (c *SMTPClient) SendEmail(toEmails []string, message []byte) error {

	err := smtp.SendMail(fmt.Sprintf(`%s:%d`, c.Server, c.Port), *c.Auth, c.Login, toEmails, message)
	if err != nil {
		return err
	}

	return nil
}
