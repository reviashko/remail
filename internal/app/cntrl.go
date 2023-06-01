package app

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/reviashko/remail/model"
	"golang.org/x/text/encoding/charmap"
)

// ControllerInterface interface
type ControllerInterface interface {
	GetUnreadMessages() ([]model.MessageInfo, error)
	GetDataForSend(msg *model.MessageInfo) (string, bool, []string)
	GenerateSubject(firstLine string) (string, bool)
	Quit()
	PrintStat()
}

// Controller struct
type Controller struct {
	Config        *model.RemailConfig
	EmailReceiver POP3ClientInterface
	EmailSender   SMTPClientInterface
	SubjectRegex  *regexp.Regexp
	DynamicParams model.DynamicParamsInterface
}

// NewController func
func NewController(emailReceiver POP3ClientInterface, config *model.RemailConfig, emailSender SMTPClientInterface, dynamicParams model.DynamicParamsInterface) Controller {

	regExp, err := regexp.Compile(`AS\d Q[0-9A-Z]+`)
	if err != nil {
		log.Fatal(err.Error())
	}

	return Controller{Config: config, EmailSender: emailSender, SubjectRegex: regExp, DynamicParams: dynamicParams, EmailReceiver: emailReceiver}
}

// GenerateSubject func
func (c *Controller) GenerateSubject(firstLine string) (string, bool) {

	if strings.HasPrefix(firstLine, "AS1 Q") {
		return strings.Replace(firstLine, "AS1", c.Config.AS1, 1), true
	}

	if strings.HasPrefix(firstLine, "AS2 Q") {
		return strings.Replace(firstLine, "AS2", c.Config.AS2, 1), true
	}

	if strings.HasPrefix(firstLine, "AS3 Q") {
		return strings.Replace(firstLine, "AS3", c.Config.AS3, 1), true
	}

	return "", false
}

// GetDataForSend func
func (c *Controller) GetDataForSend(msg *model.MessageInfo) (string, bool, []string) {

	if strings.Contains(msg.From, c.Config.OutputSource) && strings.Contains(msg.To, c.Config.OutputSource) {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(msg.Body)))
		scanner.Scan()

		subj, finded := c.GenerateSubject(c.SubjectRegex.FindString(strings.ToUpper(scanner.Text())))

		return subj, finded, []string{c.Config.OutputForward}
		//return subj, finded, []string{"devers@inbox.ru"}
	}

	return "", false, []string{}
}

// Run func
func (c *Controller) Run() {

	for true {

		time.Sleep(time.Duration(c.Config.LoopDelaySec) * time.Second)

		msgs, err := c.EmailReceiver.GetUnreadMessages(c.DynamicParams.GetLastMsgID())
		if err != nil {
			fmt.Println("GetUnreadMessages", err.Error())
			continue
		}

		if len(msgs) == 0 {
			continue
		}

		for _, m := range msgs {

			if m.MsgID > c.DynamicParams.GetLastMsgID() {
				c.DynamicParams.SetLastMsgID(m.MsgID)
			}

			subject, isSutable, to := c.GetDataForSend(&m)
			if !isSutable {
				continue
			}

			dec := charmap.Windows1251.NewEncoder()

			newBody := make([]byte, len(subject)*2)
			n, _, err2 := dec.Transform(newBody, []byte(subject), false)
			if err2 != nil {
				panic(err2)
			}
			subject = fmt.Sprintf("%s", newBody[:n])

			newBody2 := make([]byte, len(m.Body)*2)
			n2, _, err2 := dec.Transform(newBody2, []byte(m.Body), false)
			if err2 != nil {
				panic(err2)
			}
			mm := fmt.Sprintf("%s", newBody2[:n2])

			fmt.Println("id=", m.MsgID, "Subject=", subject, "from=", m.From, "to=", to)

			msg := append([]byte(fmt.Sprintf("Subject: %s\n", subject)+c.Config.WIN1251Header), mm...)

			//fmt.Println("=====> Subject=", subject, "isUTF8=", utf8.Valid([]byte(fmt.Sprintf("Subject: %s\n", subject)+c.Config.MIMEHeader)))
			fmt.Println("=====> msg=", string(msg))

			err := c.EmailSender.SendEmail(to, msg)
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		c.DynamicParams.Save()
	}
}
