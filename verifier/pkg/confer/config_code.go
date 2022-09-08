package confer

import (
	"strconv"
	"sync"
)

var (
	codeConfig sync.Map
)

type ConfigCodeList struct {
	Codes map[string]string
}

func ConfigCodeInit() {
	for k, v := range globalConfig.Code {
		codeConfig.Store(k, v)
	}
}

func ConfigCodeGetMessage(code int) string {
	msg, exists := codeConfig.Load(strconv.Itoa(code))
	if !exists {
		return "system error"
	}
	return msg.(string)
}
