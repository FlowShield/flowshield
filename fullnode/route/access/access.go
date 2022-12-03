package access

import (
	"net/http"

	v1 "github.com/flowshield/flowshield/fullnode/app/v1/access/controller"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/middle"
	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/gin-gonic/gin"
)

func checkDaoMaster(c *gin.Context) {
	// 判断当前用户是否是Dao主
	user := util.User(c)
	if user == nil || !user.Master {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    pconst.CODE_COMMON_ACCESS_FORBIDDEN,
			"message": "permission denied",
		})
		c.Abort()
	}
}

func APIAccess(parentRoute gin.IRouter) {
	r := parentRoute.Group("access", middle.Oauth2())
	//r := parentRoute.Group("access")
	{
		resource := r.Group("resource")
		{
			resource.GET("", v1.ResourceList)
			resource.POST("", checkDaoMaster, v1.AddResource)
			//resource.PUT("", v1.EditResource)
			resource.DELETE("/:uuid", checkDaoMaster, v1.DelResource)
		}
		//relay := r.Group("relay")
		//{
		//	relay.GET("", v1.RelayList)
		//	relay.POST("", v1.AddRelay)
		//	relay.PUT("", v1.EditRelay)
		//	relay.DELETE("/:uuid", v1.DelRelay)
		//}
		//server := r.Group("server")
		//{
		//	server.GET("", v1.ServerList)
		//	server.POST("", v1.AddServer)
		//	server.PUT("", v1.EditServer)
		//	server.DELETE("/:uuid", v1.DelServer)
		//}
		client := r.Group("client")
		{
			client.GET("", v1.ClientList)
			client.POST("", v1.AddClient)
			client.POST("/notify", v1.NotifyClient)
			//client.PUT("", v1.EditClient)
			//client.DELETE("/:uuid", v1.DelClient)
		}
	}

}
