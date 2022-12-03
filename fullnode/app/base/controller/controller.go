package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
)

type Res struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"message"`
}

func BindParams(c *gin.Context, params interface{}) (b bool, code int) {
	err := c.ShouldBind(params)
	if err != nil {
		logger.Warnf(c, "params error: %v", err)
		code = pconst.CODE_COMMON_PARAMS_INCOMPLETE
		return
	}
	b = true
	return
}
