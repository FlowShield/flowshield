package user

import (
	v1 "github.com/cloudslit/cloudslit/fullnode/app/v1/user/controller"
	"github.com/cloudslit/cloudslit/fullnode/pkg/middle"
	"github.com/gin-gonic/gin"
)

func APIUser(parentRoute gin.IRouter) {
	user := parentRoute.Group("user")
	{
		user.GET("/login/:company", v1.Login)
		user.GET("/oauth2/callback/:company", v1.Oauth2Callback)
		user.GET("/detail", middle.Oauth2(), v1.UserDetail)
	}
}
