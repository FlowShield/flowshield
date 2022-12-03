package route

import (
	"github.com/flowshield/flowshield/fullnode/app/base/controller"
	v1 "github.com/flowshield/flowshield/fullnode/app/v1/controlplane/controller"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/route/access"
	"github.com/flowshield/flowshield/fullnode/route/controlplane"
	"github.com/flowshield/flowshield/fullnode/route/node"
	"github.com/flowshield/flowshield/fullnode/route/user"

	"github.com/gin-gonic/gin"
)

// Home 主页
func Home(engine *gin.Engine) {
	engine.GET("", controller.Welcome)
}

func Api(engine *gin.Engine) {
	engine.GET("/a/:hash", v1.MachineOauth)
	prefix := confer.ConfigAppGetString("UrlPrefix", "")
	RouteV1 := engine.Group(prefix + pconst.APIAPIV1URL)
	{
		access.APIAccess(RouteV1)
		node.APINode(RouteV1)
		controlplane.APIControlPlane(RouteV1)
		user.APIUser(RouteV1)
	}
}

func NotFound(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "404 Not Found")
	})
}
