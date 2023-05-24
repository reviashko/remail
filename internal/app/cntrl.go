package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/reviashko/remail/model"
)

// ControllerInterface interface
type ControllerInterface interface {
	GetUnreadMessages() ([]model.MessageInfo, error)
	MakeSubjectIfSutable(msg *model.MessageInfo) (string, bool)
	Quit()
	PrintStat()
}

// Controller struct
type Controller struct {
	Config        *model.RemailConfig
	EmailReceiver POP3ClientInterface
	EmailSender   SMTPClientInterface
}

// NewController func
func NewController(emailReceiver POP3ClientInterface, config *model.RemailConfig, emailSender SMTPClientInterface) Controller {
	return Controller{EmailReceiver: emailReceiver, Config: config, EmailSender: emailSender}
}

// MakeSubjectIfSutable func
func (c *Controller) MakeSubjectIfSutable(msg *model.MessageInfo) (string, bool) {
	if msg.IsMultiPart {
		return "", false
	}

	if strings.Contains(msg.From, "<FSM-OUTSOURCE@megafon.ru>") {
		return msg.Subject, true
	}

	if strings.Contains(msg.From, "<sd@direct-credit.ru>") {
		//msg := append([]byte(fmt.Sprintf("Subject: %s\n", m.Subject)+c.Config.MIMEHeader), m.Body...)
		return msg.Subject, false
	}

	return "", false
}

// Run func
func (c *Controller) Run() {

	to := []string{
		//"devers@inbox.ru",
		"3ce2744b695b43d3a6c82f7ea0c1ff5b@webim-mail.ru",
	}

	for true {

		time.Sleep(time.Duration(c.Config.LoopDelaySec) * time.Second)

		msgs, err := c.EmailReceiver.GetUnreadMessages()
		if err != nil {
			fmt.Println(err.Error())
		}

		if len(msgs) == 0 {
			continue
		}

		for _, m := range msgs {

			subject, isSutable := c.MakeSubjectIfSutable(&m)
			if !isSutable {
				continue
			}

			fmt.Println("id=", m.MsgID, "Subject=", m.Subject, "from=", m.From)

			msg := append([]byte(fmt.Sprintf("Subject: %s\n", subject)+c.Config.MIMEHeader), m.Body...)
			err := c.EmailSender.SendEmail(to, msg)
			if err != nil {
				fmt.Println(err.Error())
			}

			/*
				//Case for resending from megafone to webim
				if strings.Contains(m.From, "<FSM-OUTSOURCE@megafon.ru>") && !m.IsMultiPart {

					fmt.Println("id=", m.MsgID, "Subject=", m.Subject, "from=", m.From)

					msg := append([]byte(fmt.Sprintf("Subject: %s\n", m.Subject)+c.Config.MIMEHeader), m.Body...)
					err := c.EmailSender.SendEmail(to, msg)
					if err != nil {
						fmt.Println(err.Error())
					}
				}

				//Case for sending from DC to megafone
				if strings.Contains(m.From, "<sd@direct-credit.ru>") && !m.IsMultiPart {
					fmt.Println("id=", m.MsgID, "Subject=", m.Subject, "from=", m.From)
				}
			*/

		}
	}
}
