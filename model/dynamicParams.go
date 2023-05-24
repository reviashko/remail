package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// DynamicParamsInterface interface
type DynamicParamsInterface interface {
	GetLastMsgID() int
	SetLastMsgID(int)
	Save() error
	InitParams() error
}

// DynamicParams struct
type DynamicParams struct {
	MsgID      int
	NeedToSave bool
}

// GetLastMsgID func
func (d *DynamicParams) GetLastMsgID() int {
	return d.MsgID
}

// SetLastMsgID func
func (d *DynamicParams) SetLastMsgID(msgID int) {
	d.MsgID = msgID
	d.NeedToSave = true
}

// Save func
func (d *DynamicParams) Save() error {
	if !d.NeedToSave {
		return nil
	}

	d.NeedToSave = false

	fileData, err := json.Marshal(d)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("dynParam.json", fileData, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Get func
func (d *DynamicParams) InitParams() error {

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadFile(fmt.Sprintf("%s/dynParam.json", pwd))
	if err != nil {
		return err
	}

	target := DynamicParams{}
	err = json.Unmarshal(fileData, &target)
	if err != nil {
		return err
	}

	d.MsgID = target.MsgID
	d.NeedToSave = false

	return nil
}
