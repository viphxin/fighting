package conf

import (
	"encoding/json"
	"github.com/viphxin/xingo/logger"
	"io/ioutil"
)

var ServerConfObj struct {
	TcpPort int
	MaxConn int
	//log
	LogPath        string
	LogName        string
	MaxLogNum      int32
	MaxFileSize    int64
	LogFileUnit    logger.UNIT
	LogLevel       logger.LEVEL
	SetToConsole   bool
	PoolSize       int32
	MaxWorkerLen   int32
	MaxSendChanLen int32
	FrameSpeed     uint8
	StepPerMs      int
}

func init() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		logger.Fatal(err)
	}
	err = json.Unmarshal(data, &ServerConfObj)
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("load conf successful!!!")
	}
}
