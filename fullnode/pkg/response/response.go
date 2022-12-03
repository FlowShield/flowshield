package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/gin-gonic/gin"
)

func UtilResponseReturnJson(c *gin.Context, code int, model interface{}, msg ...string) {
	UtilResponseReturnJsonWithMsg(c, code, GetResponseMsg(c, code), model, true, true)
}

func UtilResponseReturnJsonWithMsg(c *gin.Context, code int, msg string, model interface{},
	callbackFlag bool, unifyCode bool) {
	if unifyCode && code == 0 {
		code = 1001
	}
	//if msg == "" {
	//	msg = confer.ConfigCodeGetMessage(code)
	//}
	// 放入返回的code码
	//c.Set("result_code", code)
	//// 判断是否存在error上下文
	//if err, ok := c.Get("error"); ok {
	//	if err.(error) != nil {
	//		msg = err.(error).Error()
	//	}
	//}
	rj := gin.H{
		"code":    code,
		"message": msg,
		"data":    model,
	}
	var callback string
	if callbackFlag {
		callback = c.Query("callback")
	}
	if len(strings.TrimSpace(callback)) == 0 {
		if httpStatus := c.GetInt("httpStatus"); httpStatus > 0 {
			c.JSON(httpStatus, rj)
		} else if status, ok := statusCode[code]; ok {
			// 根据code码返回不同的statusCode
			c.JSON(status, rj)
		} else {
			c.JSON(http.StatusOK, rj)
		}
	} else {
		if r, err := json.Marshal(rj); err == nil {
			c.String(http.StatusOK, "%s(%s)", callback, r)
		}
	}
}

func UtilResponseReturnJsonFailed(c *gin.Context, code int) {
	UtilResponseReturnJson(c, code, nil)
}

func GetResponseCode(code int) int {
	return statusCode[code]
}

func GetResponseMsg(gin *gin.Context, code int) (message string) {
	if message = gin.GetString("error"); len(message) > 0 {
		return
	}
	var lang string
	if gin != nil {
		lang = util.GetCookieFromGin(gin, "lang")
	}
	if len(lang) == 0 {
		lang = "zh-CN"
	}
	switch lang {
	case "zh-CN", "en-US":
		return messageText[lang][code]
	default:
		return messageText["zh-CN"][code]
	}
}

// 创建错误码和statusCode的msp关系
var statusCode = map[int]int{
	pconst.CODE_ERROR_OK:                  http.StatusOK,
	pconst.CODE_COMMON_OK:                 http.StatusOK,
	pconst.CODE_COMMON_SERVER_BUSY:        http.StatusInternalServerError,
	pconst.CODE_COMMON_PARAMS_INCOMPLETE:  http.StatusBadRequest,
	pconst.CODE_COMMON_DATA_NOT_EXIST:     http.StatusBadRequest,
	pconst.CODE_COMMON_DATA_ALREADY_EXIST: http.StatusBadRequest,
	pconst.CODE_DATA_HAS_RELATION:         http.StatusBadRequest,
	pconst.CODE_DATA_WRONG_STATE:          http.StatusBadRequest,
	pconst.CODE_API_BAD_REQUEST:           http.StatusBadRequest,
	pconst.CODE_VICTORIA_METRICS_ERR:      http.StatusInternalServerError,
	pconst.CODE_NOT_FOUND:                 http.StatusNotFound,
}

var messageText = map[string]map[int]string{
	"zh-CN": messageCNText,
	"en-US": messageENText,
}

var messageCNText = map[int]string{
	pconst.CODE_ERROR_OK:                  "Success",
	pconst.CODE_COMMON_OK:                 "Success",
	pconst.CODE_COMMON_SERVER_BUSY:        "服务器繁忙，请稍后再试",
	pconst.CODE_COMMON_PARAMS_INCOMPLETE:  "请求参数错误",
	pconst.CODE_COMMON_DATA_NOT_EXIST:     "记录不存在",
	pconst.CODE_COMMON_DATA_ALREADY_EXIST: "记录已经存在",
	pconst.CODE_DATA_HAS_RELATION:         "存在关联数据，无法删除",
	pconst.CODE_DATA_WRONG_STATE:          "数据状态出错",
	pconst.CODE_API_BAD_REQUEST:           "接口请求失败",
	pconst.CODE_VICTORIA_METRICS_ERR:      "数据出错，请稍后重试",
	pconst.CODE_NOT_FOUND:                 "资源不存在",
}

var messageENText = map[int]string{
	pconst.CODE_ERROR_OK:                  "Success",
	pconst.CODE_COMMON_OK:                 "Success",
	pconst.CODE_COMMON_SERVER_BUSY:        "Internal Server Error,Try Again Later",
	pconst.CODE_COMMON_PARAMS_INCOMPLETE:  "Parameter Error",
	pconst.CODE_COMMON_DATA_NOT_EXIST:     "Record Does Not Exist",
	pconst.CODE_COMMON_DATA_ALREADY_EXIST: "Record Already Exists",
	pconst.CODE_DATA_HAS_RELATION:         "The Data Is Associated,Can Not Be Deleted",
	pconst.CODE_DATA_WRONG_STATE:          "Data State Is Wrong",
	pconst.CODE_API_BAD_REQUEST:           "API Request Error",
	pconst.CODE_VICTORIA_METRICS_ERR:      "Data Error,Try Again Later",
	pconst.CODE_NOT_FOUND:                 "Resource Not Found",
}
