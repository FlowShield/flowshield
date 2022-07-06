package node

import (
	v1 "github.com/cloudslit/cloudslit/fullnode/app/v1/node/controller"
	"github.com/cloudslit/cloudslit/fullnode/pkg/middle"
	"github.com/gin-gonic/gin"
)

func APINode(parentRoute gin.IRouter) {
	node := parentRoute.Group("node", middle.Oauth2())
	{
		node.GET("", v1.ListNode)
	}
}
