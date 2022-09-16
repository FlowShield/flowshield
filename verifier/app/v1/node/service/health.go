package service

import (
	"github.com/cloudslit/cloudslit/verifier/app/v1/node/model/mparam"
	"github.com/cloudslit/cloudslit/verifier/verify"
	"github.com/gin-gonic/gin"
)

func ProviderHealth(c *gin.Context, param *mparam.ProviderHealth) (code int, data map[string]*verify.Record) {
	data = verify.VerObj.ProviderHealth(param.OrderID)
	return
}
