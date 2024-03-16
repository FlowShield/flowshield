package v1

import (
	"github.com/flowshield/flowshield/fullnode/app/base/controller"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mparam"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/service"
	"github.com/flowshield/flowshield/fullnode/pkg/response"
	"github.com/gin-gonic/gin"
)

// @Summary ClientList
// @Description Get ZTA client
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
// @Description Added ZTA client
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param Client body mparam.AddClient true "Added ZTA client"
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
//// @Description Modify ZTA client
//// @Tags ZTA
//// @Accept  json
//// @Produce  json
//// @Param Client body mparam.EditClient true "Modify ZTA client"
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
//// @Description Delete ZTA client
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
// @Description Receive notifications of client order status changes
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param NotifyClient body mparam.NotifyClient true "Receive notifications of client order status changes"
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
