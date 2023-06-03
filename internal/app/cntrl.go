package app

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"mime"
	"regexp"
	"strings"
	"time"

	"github.com/reviashko/remail/model"
)

// ControllerInterface interface
type ControllerInterface interface {
	GetUnreadMessages() ([]model.MessageInfo, error)
	GetDataForSend(msg *model.MessageInfo) (string, bool, []string)
	GenerateSubject(firstLine string) (string, bool)
	LoopFunc() error
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
	}

	return "", false, []string{}
}

// LoopFunc func
func (c *Controller) LoopFunc() error {

	msgs, err := c.EmailReceiver.GetUnreadMessages(c.DynamicParams.GetLastMsgID())
	if err != nil {
		return err
	}

	if len(msgs) == 0 {
		return nil
	}

	for _, m := range msgs {

		if m.MsgID > c.DynamicParams.GetLastMsgID() {
			c.DynamicParams.SetLastMsgID(m.MsgID)
		}

		subject, isSutable, to := c.GetDataForSend(&m)
		if !isSutable {
			continue
		}

		t := time.Now()
		fmt.Printf("[%02d-%02dT%02d:%02d] id=%d subj=[%s] from=[%s] to=[%s]\n", t.Day(), t.Month(), t.Hour(), t.Minute(), m.MsgID, subject, m.From, to)

		m.Subject = mime.QEncoding.Encode("utf-8", subject)
		err := c.EmailSender.SendEmail(to, m.ToBytes())
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	return c.DynamicParams.Save()
}

// Run func
func (c *Controller) Run() {

	for true {

		time.Sleep(time.Duration(c.Config.LoopDelaySec) * time.Second)

		err := c.LoopFunc()
		if err != nil {
			t := time.Now()
			fmt.Printf("[%02d-%02dT%02d:%02d] MainLoop=[%s]\n", t.Day(), t.Month(), t.Hour(), t.Minute(), err.Error())
		}
	}
}
