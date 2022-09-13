package server

import (
	"fmt"
	"strconv"

	"github.com/cloudslit/cloudslit/verifier/pkg/confer"
	"github.com/cloudslit/cloudslit/verifier/pkg/gin"
	"github.com/cloudslit/cloudslit/verifier/pkg/middle"
	"github.com/cloudslit/cloudslit/verifier/route"
)

func RunHTTP() {
	engine := gin.NewGin()
	// 跨域中间件
	engine.Use(middle.CorsV2())
	route.Home(engine)
	route.Api(engine)
	route.NotFound(engine)
	httpPort := confer.ConfigAppGetInt("port", 80)
	portStr := ":" + strconv.Itoa(httpPort)
	fmt.Println("start", httpPort)
	gin.ListenHTTP(portStr, engine, 10)
}
