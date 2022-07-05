package node

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/cloudslit/cloudslit/fullnode/app/v1/node/controller"
)

func APINode(parentRoute gin.IRouter) {
	node := parentRoute.Group("node")
	{
		node.GET("", v1.ListNode)
	}
}
