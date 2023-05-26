package main

import (
	"github.com/reviashko/remail/internal/app"

	"github.com/reviashko/remail/model"
)

func main() {

	appConfig := model.RemailConfig{}
	appConfig.InitParams()

	dynParams := model.DynamicParams{}
	dynParams.InitParams()

	pop3Client := app.NewPOP3Client(appConfig.POP3Host, appConfig.POP3Port, appConfig.TLSEnabled, appConfig.Login, appConfig.Pswd)
	smtpClient := app.NewSMTPClient(appConfig.SMTPHost, appConfig.SMTPPort, appConfig.Login, appConfig.Pswd)

	cntrl := app.NewController(&pop3Client, &appConfig, &smtpClient, &dynParams)
	cntrl.Run()
}
