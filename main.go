package main

import (
	"github.com/reviashko/remail/internal/app"

	"github.com/reviashko/remail/model"
)

func main() {

	appConfig := model.RemailConfig{}
	appConfig.InitParams()

	//TODO: save proceeded emails ids
	//TODO: init it here
	lastMsgID := 0

	pop3Client := app.NewPOP3Client(appConfig.POP3Host, appConfig.POP3Port, appConfig.TLSEnabled, appConfig.Login, appConfig.Pswd, lastMsgID)
	defer pop3Client.Quit()

	smtpClient := app.NewSMTPClient(appConfig.SMTPHost, appConfig.SMTPPort, appConfig.Login, appConfig.Pswd)

	cntrl := app.NewController(&pop3Client, &appConfig, &smtpClient)
	cntrl.Run()
}
