package node

import (
	v1 "github.com/flowshield/flowshield/fullnode/app/v1/node/controller"
	"github.com/flowshield/flowshield/fullnode/pkg/middle"
	"github.com/gin-gonic/gin"
)

func APINode(parentRoute gin.IRouter) {
	node := parentRoute.Group("node", middle.Oauth2())
	{
		node.GET("", v1.ListNode)
	}
}
