package logs

import (
	seelog "github.com/axengine/seelog"
)

var Logger seelog.LoggerInterface

func init() {
	Logger = seelog.Disabled
	logger, err := seelog.LoggerFromConfigAsFile("./conf/log.xml")
	if err != nil {
		panic(err)
	}
	Logger = logger
}
