package v1

import (
	"github.com/cloudslit/cloudslit/verifier/app/base/controller"
	"github.com/cloudslit/cloudslit/verifier/app/v1/node/model/mparam"
	"github.com/cloudslit/cloudslit/verifier/app/v1/node/service"
	"github.com/cloudslit/cloudslit/verifier/pkg/response"
	"github.com/gin-gonic/gin"
)

// @Summary ListNode
// @Description 获取全部节点
// @Tags ZTA Node
// @Produce  json
// @Success 200 {object} controller.Res
// @Router /verifier/provider/health [post]
func ProviderHealth(c *gin.Context) {
	param := &mparam.ProviderHealth{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code, data := service.ProviderHealth(c, param)
	response.UtilResponseReturnJson(c, code, data)
}
