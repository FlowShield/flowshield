package core

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MSPNormalStateHTTPStatusCode int64 = 1001
)

type WrapContext struct {
	Ctx context.Context
	G   *gin.Context
	Is  *I
}

type WrapHandler func(c *WrapContext) (interface{}, error)

type ErrCode string

type WrapErrorResponse struct {
	Msg     string  `json:"msg,omitempty"`
	ErrCode ErrCode `json:"errCode,omitempty"`
}

type MSPNormalizeHTTPResponseBody struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type MSPNormalizeListPaginateParams struct {
	Page     int `json:"page" form:"page"`
	LimitNum int `json:"limit_num" form:"limit_num"`
}

type MSPNormalizeList struct {
	List     interface{}          `json:"list"`
	Paginate MSPNormalizePaginate `json:"paginate"`
}

type MSPNormalizePaginate struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

func mspErrCodeTransform(x interface{}) int64 {
	switch x := x.(type) {
	case int:
		return int64(x)
	case int64:
		return x
	case string:
		i, err := strconv.Atoi(x)
		if err != nil {
			return 0
		}
		return int64(i)
	default:
		return -1
	}
}

// BindG bind params or exit
func (c *WrapContext) BindG(params interface{}) {
	err := c.G.ShouldBind(params)
	if err != nil {
		c.G.AbortWithStatusJSON(http.StatusBadRequest, WrapErrorResponse{
			Msg:     err.Error(),
			ErrCode: "1000",
		})
		panic(err)
	}
}

// WrapH ...
func WrapH(h WrapHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				Is.Logger.Warnf("request error: %s", err)
			}
		}()

		var wrapCtx WrapContext
		wrapCtx.G = c
		wrapCtx.Is = Is
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		wrapCtx.Ctx = ctx

		res, err := h(&wrapCtx)
		if err != nil {
			if errCode, ok := res.(ErrCode); ok {
				errCodeString := string(errCode)
				if len(errCodeString) > 3 {
					c.JSON(http.StatusBadRequest, MSPNormalizeHTTPResponseBody{
						Code:    mspErrCodeTransform(errCodeString),
						Message: err.Error(),
					})
					return
				}
				errCodeInt, _ := strconv.Atoi(string(errCode))
				c.JSON(errCodeInt, MSPNormalizeHTTPResponseBody{
					Code:    mspErrCodeTransform(errCodeInt),
					Message: err.Error(),
				})
				return
			}
			if intRes, ok := res.(int); ok {
				c.JSON(intRes, MSPNormalizeHTTPResponseBody{
					Code:    mspErrCodeTransform(intRes),
					Message: err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, MSPNormalizeHTTPResponseBody{
				Code:    mspErrCodeTransform(500),
				Message: err.Error(),
			})
			return
		}

		if intRes, ok := res.(int); ok {
			if len(strconv.Itoa(intRes)) == 3 {
				c.JSON(intRes, MSPNormalizeHTTPResponseBody{
					Code:    MSPNormalStateHTTPStatusCode,
					Message: "Success!",
				})
				return
			}
		}

		if res == nil {
			c.JSON(http.StatusNoContent, MSPNormalizeHTTPResponseBody{
				Code:    MSPNormalStateHTTPStatusCode,
				Message: "No Content.",
			})
			return
		}

		c.JSON(http.StatusOK, MSPNormalizeHTTPResponseBody{
			Code:    MSPNormalStateHTTPStatusCode,
			Data:    res,
			Message: "Success!",
		})
	}
}
