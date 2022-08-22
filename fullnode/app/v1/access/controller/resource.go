package v1

import (
	"strings"

	"github.com/cloudslit/cloudslit/fullnode/app/base/controller"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/service"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	"github.com/cloudslit/cloudslit/fullnode/pkg/response"
	"github.com/cloudslit/cloudslit/fullnode/pkg/util"

	"github.com/gin-gonic/gin"
)

// @Summary ResourceList
// @Description 获取ZTA的resource
// @Tags ZTA
// @Produce  json
// @Success 200 {object} controller.Res
// @Router /access/resource [get]
func ResourceList(c *gin.Context) {
	param := mparam.ResourceList{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code, data := service.ResourceList(c, param)
	response.UtilResponseReturnJson(c, code, data)
}

// @Summary AddResource
// @Description 新增ZTA的resource
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param Resource body mparam.AddResource true "新增ZTA的resource"
// @Success 200 {object} controller.Res
// @Router /access/resource [post]
func AddResource(c *gin.Context) {
	param := &mparam.AddResource{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	if len(param.Type) > 0 && param.Type == "cidr" {
		// 判断是不是纯IP格式
		if strings.Contains(param.Host, "/") {
			if !util.IsCIDR(param.Host) {
				response.UtilResponseReturnJsonFailed(c, pconst.CODE_COMMON_PARAMS_INCOMPLETE)
				return
			}
		} else {
			if !util.IsIP(param.Host) {
				response.UtilResponseReturnJsonFailed(c, pconst.CODE_COMMON_PARAMS_INCOMPLETE)
				return
			}
		}
	}
	code = service.AddResource(c, param)
	response.UtilResponseReturnJson(c, code, nil)
}

//// @Summary EditResource
//// @Description 修改ZTA的resource
//// @Tags ZTA
//// @Accept  json
//// @Produce  json
//// @Param Resource body mparam.EditResource true "修改ZTA的resource"
//// @Success 200 {object} controller.Res
//// @Router /access/resource [put]
//func EditResource(c *gin.Context) {
//	param := &mparam.EditResource{}
//	b, code := controller.BindParams(c, &param)
//	if !b {
//		response.UtilResponseReturnJsonFailed(c, code)
//		return
//	}
//	if len(param.Type) > 0 && param.Type == "cidr" {
//		// 判断是不是纯IP格式
//		if strings.Contains(param.Host, "/") {
//			if !util.IsCIDR(param.Host) {
//				response.UtilResponseReturnJsonFailed(c, pconst.CODE_COMMON_PARAMS_INCOMPLETE)
//				return
//			}
//		} else {
//			if !util.IsIP(param.Host) {
//				response.UtilResponseReturnJsonFailed(c, pconst.CODE_COMMON_PARAMS_INCOMPLETE)
//				return
//			}
//		}
//	}
//	code = service.EditResource(c, param)
//	response.UtilResponseReturnJson(c, code, nil)
//}

// @Summary DelResource
// @Description 删除ZTA的resource
// @Tags ZTA
// @Produce  json
// @Param id path int true "主键ID"
// @Success 200 {object} controller.Res
// @Router /access/resource/{uuid} [delete]
func DelResource(c *gin.Context) {
	code := service.DelResource(c, c.Param("uuid"))
	response.UtilResponseReturnJson(c, code, nil)
}
