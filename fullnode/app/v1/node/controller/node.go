package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/flowshield/flowshield/fullnode/app/base/controller"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mparam"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/service"
	"github.com/flowshield/flowshield/fullnode/pkg/response"
)

// @Summary ListNode
// @Description 获取全部节点
// @Tags ZTA Node
// @Produce  json
// @Success 200 {object} controller.Res
// @Router /node [get]
func ListNode(c *gin.Context) {
	param := mparam.ListNode{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code, data := service.ListNode(c, param)
	response.UtilResponseReturnJson(c, code, data)
}

//
//// @Summary AddOauth2
//// @Description 新增ZTA的Oauth2
//// @Tags ZTA Oauth2
//// @Accept  json
//// @Produce  json
//// @Param Oauth2 body mparam.AddOauth2 true "新增ZTA的Oauth2"
//// @Success 200 {object} controller.Res
//// @Router /sysytem/oauth2 [post]
//func AddOauth2(c *gin.Context) {
//	param := &mmysql.Oauth2{}
//	b, code := controller.BindParams(c, &param)
//	if !b {
//		response.UtilResponseReturnJsonFailed(c, code)
//		return
//	}
//	code = service.AddOauth2(c, param)
//	response.UtilResponseReturnJson(c, code, nil)
//}
//
//// @Summary EditOauth2
//// @Description 修改ZTA的Oauth2
//// @Tags ZTA Oauth2
//// @Accept  json
//// @Produce  json
//// @Param Oauth2 body mparam.EditOauth2 true "修改ZTA的Oauth2"
//// @Success 200 {object} controller.Res
//// @Router /sysytem/oauth2 [put]
//func EditOauth2(c *gin.Context) {
//	param := &mparam.EditOauth2{}
//	b, code := controller.BindParams(c, &param)
//	if !b {
//		response.UtilResponseReturnJsonFailed(c, code)
//		return
//	}
//	code = service.EditOauth2(c, param)
//	response.UtilResponseReturnJson(c, code, nil)
//}
//
//// @Summary DelOauth2
//// @Description 删除ZTA的Oauth2
//// @Tags ZTA Oauth2
//// @Produce  json
//// @Param id path int true "主键ID"
//// @Success 200 {object} controller.Res
//// @Router /sysytem/oauth2/{id} [delete]
//func DelOauth2(c *gin.Context) {
//	id := c.Param("id")
//	idInt, err := strconv.Atoi(id)
//	if err != nil {
//		response.UtilResponseReturnJsonFailed(c, pconst.CODE_COMMON_PARAMS_INCOMPLETE)
//		return
//	}
//	code := service.DelOauth2(c, uint64(idInt))
//	response.UtilResponseReturnJson(c, code, nil)
//}
