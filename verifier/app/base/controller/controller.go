package controller

import (
	"github.com/flowshield/flowshield/verifier/pconst"
	"github.com/flowshield/flowshield/verifier/pkg/logger"
	"github.com/gin-gonic/gin"
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
