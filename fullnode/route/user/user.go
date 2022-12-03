package user

import (
	v1 "github.com/flowshield/flowshield/fullnode/app/v1/user/controller"
	"github.com/flowshield/flowshield/fullnode/pkg/middle"
	"github.com/gin-gonic/gin"
)

func APIUser(parentRoute gin.IRouter) {
	user := parentRoute.Group("user")
	{
		user.GET("/login", v1.Login)
		user.GET("/oauth2/callback", v1.Oauth2Callback)
		user.GET("/detail", middle.Oauth2(), v1.UserDetail)
		user.POST("/refresh", middle.Oauth2(), v1.UserRefresh)
	}
}
