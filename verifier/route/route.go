package route

import (
	"github.com/flowshield/flowshield/verifier/app/base/controller"
	"github.com/flowshield/flowshield/verifier/pconst"
	"github.com/flowshield/flowshield/verifier/pkg/confer"
	"github.com/flowshield/flowshield/verifier/route/health"
	"github.com/gin-gonic/gin"
)

// Home 主页
func Home(engine *gin.Engine) {
	engine.GET("", controller.Welcome)
}

func Api(engine *gin.Engine) {
	prefix := confer.ConfigAppGetString("UrlPrefix", "")
	RouteV1 := engine.Group(prefix + pconst.APIAPIV1URL)
	{
		health.APIHealth(RouteV1)
	}
}

func NotFound(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "404 Not Found")
	})
}
