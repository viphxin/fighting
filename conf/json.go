package conf

import (
	"encoding/json"
	"github.com/viphxin/xingo/logger"
	"io/ioutil"
)

var ServerConfObj struct {
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