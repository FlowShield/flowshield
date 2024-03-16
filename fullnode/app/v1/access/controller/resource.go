package v1

import (
	"strings"

	"github.com/flowshield/flowshield/fullnode/app/base/controller"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mparam"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/service"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/response"
	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/gin-gonic/gin"
)

// @Summary ResourceList
// @Description Get ZTA resources
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
// @Description Added ZTA resources
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param Resource body mparam.AddResource true "Added ZTA resources"
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
		// Determine whether it is pure IP format
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
//// @Description Modify ZTA resources
//// @Tags ZTA
//// @Accept  json
//// @Produce  json
//// @Param Resource body mparam.EditResource true "Modify ZTA resources"
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
//		// Determine whether it is pure IP format
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
// @Description Delete ZTA resources
// @Tags ZTA
// @Produce  json
// @Param id path int true "Primary key ID"
// @Success 200 {object} controller.Res
// @Router /access/resource/{uuid} [delete]
func DelResource(c *gin.Context) {
	code := service.DelResource(c, c.Param("uuid"))
	response.UtilResponseReturnJson(c, code, nil)
}
