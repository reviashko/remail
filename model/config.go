package model

import (
	"log"

	"github.com/tkanos/gonfig"
)

// RemailConfig struct
type RemailConfig struct {
	POP3Host      string
	POP3Port      int
	SMTPHost      string
	SMTPPort      int
	TLSEnabled    bool
	Login         string
	Pswd          string
	LoopDelaySec  int
	AS1           string
	AS2           string
	AS3           string
	OutputForward string
	OutputSource  string
}

// InitParams func
func (c *RemailConfig) InitParams() {

	if err := gonfig.GetConf("config/data.json", c); err != nil {
		log.Panicf("load spec confg error: %s\n", err.Error())
	}
}
