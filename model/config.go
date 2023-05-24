package model

import (
	"log"
	"strings"

	"github.com/tkanos/gonfig"
)

// RemailConfig struct
type RemailConfig struct {
	POP3Host     string
	POP3Port     int
	SMTPHost     string
	SMTPPort     int
	TLSEnabled   bool
	Login        string
	Pswd         string
	LoopDelaySec int
	Patterns     []string
	MIMEHeader   string
}

// InitParams func
func (c *RemailConfig) InitParams() {

	if err := gonfig.GetConf("config/data.json", c); err != nil {
		log.Panicf("load spec confg error: %s\n", err.Error())
	}
}

// IsPatternMatched func
func (c *RemailConfig) IsPatternMatched(data string) bool {

	for _, p := range c.Patterns {
		if strings.Contains(data, p) {
			return true
		}
	}

	return false
}
