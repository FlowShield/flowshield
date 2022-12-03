package controller

import (
	"github.com/flowshield/flowshield/fullnode/app/base/controller"
	"github.com/flowshield/flowshield/fullnode/app/v1/controlplane/model/mparam"
	"github.com/flowshield/flowshield/fullnode/app/v1/controlplane/service"
	"github.com/flowshield/flowshield/fullnode/pkg/response"
	"github.com/gin-gonic/gin"
)

// @Summary LoginUrl
// @Description 根据机器码获取客户端鉴权的url
// @Tags ZTA ControlPlane
// @Produce  json
// @Param machine_id path string true "machine_id"
// @Success 200 {object} controller.Res
// @Router /controlplane/machine/{machine_id} [get]
func LoginUrl(c *gin.Context) {
	code, data := service.GetLoginUrl(c, c.Param("machine_id"))
	response.UtilResponseReturnJson(c, code, data)
}

// @Summary MachineOauth
// @Description 机器鉴权
// @Tags ZTA ControlPlane
// @Produce  json
// @Param hash path string true "hash"
// @Success 200 {object} controller.Res
// @Router /a/{hash} [get]
func MachineOauth(c *gin.Context) {
	service.MachineOauth(c, c.Param("hash"))
}

// @Summary MachineLongPoll
// @Description 机器鉴权
// @Tags ZTA ControlPlane
// @Produce  json
// @Param category query string true "轮询的主题"
// @Param timeout query int true "超时时间，单位：秒"
// @Success 200 {object} controller.Res
// @Router /machine/auth/poll [get]
func MachineLongPoll(c *gin.Context) {
	param := mparam.MachineLongPoll{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	data, code := service.MachineLongPoll(c, param)
	response.UtilResponseReturnJson(c, code, data)
}
