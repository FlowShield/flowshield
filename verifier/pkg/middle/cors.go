package middle

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 用于跨域头维护
var corsHeader = []string{""}

// CorsV2 跨域
func CorsV2() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders(corsHeader...)
	return cors.New(config)
}
