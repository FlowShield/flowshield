package service

import (
	"github.com/flowshield/flowshield/verifier/app/v1/node/model/mparam"
	"github.com/flowshield/flowshield/verifier/verify"
	"github.com/gin-gonic/gin"
)

func ProviderHealth(c *gin.Context, param *mparam.ProviderHealth) (code int, data map[string]*verify.Record) {
	data = verify.VerObj.ProviderHealth(param.OrderID)
	return
}
