package main

import (
	"fmt"
	"testing"

	"github.com/reviashko/remail/internal/app"
	"github.com/reviashko/remail/model"
)

// POP3MockClient struct
type POP3MockClient struct {
	LastMsgID int
}

// GetUnreadMessages func
func (p *POP3MockClient) GetUnreadMessages(msgID int) ([]model.MessageInfo, error) {

	retval := make([]model.MessageInfo, 0)

	m := 1
	for m <= p.LastMsgID {
		retval = append(retval, model.MessageInfo{MsgID: m, To: fmt.Sprintf("mail%d@mail.ru", m), From: fmt.Sprintf("mail%d_%d@mail.ru", m, m), Body: []byte(fmt.Sprintf("AS%d Q%d%d%d%d%d%d%d\nmessage body number %d", m, m, m, m, m, m, m, m, m))})
		m++
	}

	return retval, nil
}

// SMTPMockClient struct
type SMTPMockClient struct{}

// SendEmail func
func (c *SMTPMockClient) SendEmail(toEmails []string, message []byte) error {

	fmt.Printf("SendMail to: [%v]\nBody:[%s]", toEmails, message)
	return nil
}

// DynamicMockParams struct
type DynamicMockParams struct {
	MsgID      int
	NeedToSave bool
}

// GetLastMsgID func
func (d *DynamicMockParams) GetLastMsgID() int {
	return d.MsgID
}

// SetLastMsgID func
func (d *DynamicMockParams) SetLastMsgID(msgID int) {
	d.MsgID = msgID
	d.NeedToSave = true
}

// Save func
func (d *DynamicMockParams) Save() error {
	if !d.NeedToSave {
		return nil
	}

	d.NeedToSave = false
	return nil
}

// Get func
func (d *DynamicMockParams) InitParams() error {
	return nil
}

func TestController(t *testing.T) {

	appConfig := model.RemailConfig{}
	appConfig.InitParams()

	dynParams := DynamicMockParams{}
	dynParams.InitParams()

	pop3Client := POP3MockClient{LastMsgID: 6}
	smtpClient := SMTPMockClient{}

	cntrl := app.NewController(&pop3Client, &appConfig, &smtpClient, &dynParams)

	// first case
	// get 6 email messages and controller dont going to react it, because dynParams last MsgID is 10 (> 6)
	test1MaxMsgID := 10
	dynParams.SetLastMsgID(test1MaxMsgID)
	dynParams.Save()

	err := cntrl.LoopFunc()
	if err != nil {
		t.Errorf("LoopFunc: %s", err.Error())
	}

	if dynParams.MsgID > test1MaxMsgID {
		t.Errorf("MsgID %d > %d, wait for %d", dynParams.MsgID, test1MaxMsgID, test1MaxMsgID)
	}

	// second case
	// get 6 email messages and controller must going to react it, because dynParams last MsgID is 5 (< 6)
	tes21MaxMsgID := 5
	dynParams.SetLastMsgID(tes21MaxMsgID)
	dynParams.Save()

	err = cntrl.LoopFunc()
	if err != nil {
		t.Errorf("LoopFunc: %s", err.Error())
	}

	if dynParams.MsgID < tes21MaxMsgID {
		t.Errorf("MsgID %d < %d, wait for %d", dynParams.MsgID, tes21MaxMsgID, tes21MaxMsgID)
	}
}
