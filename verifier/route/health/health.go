package health

import (
	v1 "github.com/flowshield/flowshield/verifier/app/v1/node/controller"
	"github.com/gin-gonic/gin"
)

func APIHealth(parentRoute gin.IRouter) {
	node := parentRoute.Group("verifier/provider/health")
	{
		node.POST("", v1.ProviderHealth)
	}
}
