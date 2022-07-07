package v1

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/controller"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/service"
	"github.com/cloudslit/cloudslit/fullnode/pkg/response"
	"github.com/gin-gonic/gin"
)

// @Summary ClientList
// @Description 获取ZTA的client
// @Tags ZTA
// @Produce  json
// @Success 200 {object} controller.Res
// @Router /access/client [get]
func ClientList(c *gin.Context) {
	param := mparam.ClientList{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code, data := service.ClientList(c, param)
	response.UtilResponseReturnJson(c, code, data)
}

// @Summary AddClient
// @Description 新增ZTA的client
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param Client body mparam.AddClient true "新增ZTA的client"
// @Success 200 {object} controller.Res
// @Router /access/client [post]
func AddClient(c *gin.Context) {
	param := &mparam.AddClient{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code, data := service.AddClient(c, param)
	response.UtilResponseReturnJson(c, code, data)
}

//// @Summary EditClient
//// @Description 修改ZTA的client
//// @Tags ZTA
//// @Accept  json
//// @Produce  json
//// @Param Client body mparam.EditClient true "修改ZTA的client"
//// @Success 200 {object} controller.Res
//// @Router /access/client [put]
//func EditClient(c *gin.Context) {
//	param := &mparam.EditClient{}
//	b, code := controller.BindParams(c, &param)
//	if !b {
//		response.UtilResponseReturnJsonFailed(c, code)
//		return
//	}
//	code = service.EditClient(c, param)
//	response.UtilResponseReturnJson(c, code, nil)
//}
//
//// @Summary DelClient
//// @Description 删除ZTA的client
//// @Tags ZTA
//// @Produce  json
//// @Param uuid path string true "uuid"
//// @Success 200 {object} controller.Res
//// @Router /access/client/{uuid} [delete]
//func DelClient(c *gin.Context) {
//	code := service.DelClient(c, c.Param("uuid"))
//	response.UtilResponseReturnJson(c, code, nil)
//}

// @Summary NotifyClient
// @Description 接收client订单状态改变的通知
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param NotifyClient body mparam.NotifyClient true "接收client订单状态改变的通知"
// @Success 200 {object} controller.Res
// @Router /access/client/notify [post]
func NotifyClient(c *gin.Context) {
	param := &mparam.NotifyClient{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code = service.NotifyClient(c, param)
	response.UtilResponseReturnJson(c, code, nil)
}
