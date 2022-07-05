package access

import (
	v1 "github.com/cloudslit/cloudslit/fullnode/app/v1/access/controller"
	"github.com/cloudslit/cloudslit/fullnode/pkg/middle"

	"github.com/gin-gonic/gin"
)

func APIAccess(parentRoute gin.IRouter) {
	r := parentRoute.Group("access", middle.Oauth2())
	//r := parentRoute.Group("access")
	{
		resource := r.Group("resource")
		{
			resource.GET("", v1.ResourceList)
			resource.POST("", v1.AddResource)
			resource.PUT("", v1.EditResource)
			resource.DELETE("/:uuid", v1.DelResource)
		}
		relay := r.Group("relay")
		{
			relay.GET("", v1.RelayList)
			relay.POST("", v1.AddRelay)
			relay.PUT("", v1.EditRelay)
			relay.DELETE("/:uuid", v1.DelRelay)
		}
		server := r.Group("server")
		{
			server.GET("", v1.ServerList)
			server.POST("", v1.AddServer)
			server.PUT("", v1.EditServer)
			server.DELETE("/:uuid", v1.DelServer)
		}
		client := r.Group("client")
		{
			client.GET("", v1.ClientList)
			client.POST("", v1.AddClient)
			client.PUT("", v1.EditClient)
			client.DELETE("/:uuid", v1.DelClient)
		}
	}

}
