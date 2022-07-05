package helper

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	// MSPNormalStateHTTPStatusCode ...
	MSPNormalStateHTTPStatusCode int64 = 1001
)

// HTTPWrapContext ...
type HTTPWrapContext struct {
	Ctx context.Context
	G   *gin.Context
}

// HTTPWrapHandler ...
type HTTPWrapHandler func(c *HTTPWrapContext) (interface{}, error)

// HTTPErrCode ...
type HTTPErrCode string

// HTTPWrapErrorResponse ...
type HTTPWrapErrorResponse struct {
	Msg     string      `json:"msg,omitempty"`
	ErrCode HTTPErrCode `json:"errCode,omitempty"`
}

// MSPNormalizeHTTPResponseBody ...
type MSPNormalizeHTTPResponseBody struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// DefaultMSPNormalizeListPaginateParams ...
var DefaultMSPNormalizeListPaginateParams = MSPNormalizeListPaginateParams{
	Page:     1,
	LimitNum: 20,
}

// MSPNormalizeListPaginateParams ...
type MSPNormalizeListPaginateParams struct {
	Page     int `json:"page" form:"page"`
	LimitNum int `json:"limit_num" form:"limit_num"`
}

// MSPNormalizeList ...
type MSPNormalizeList struct {
	List     interface{}          `json:"list"`
	Paginate MSPNormalizePaginate `json:"paginate"`
}

// MSPNormalizePaginate ...
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

// BindG ...
func (c *HTTPWrapContext) BindG(params interface{}) {
	err := c.G.ShouldBind(params)
	if err != nil {
		c.G.AbortWithStatusJSON(http.StatusBadRequest, HTTPWrapErrorResponse{
			Msg:     err.Error(),
			ErrCode: "1000",
		})
		panic(err)
	}
}

// WrapH ...
func WrapH(h HTTPWrapHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Named("router").Warnf("request error: %v", err)
			}
		}()

		var wrapCtx HTTPWrapContext
		wrapCtx.G = c
		var cancel context.CancelFunc
		wrapCtx.Ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := h(&wrapCtx)
		if err != nil {
			if errCode, ok := res.(HTTPErrCode); ok {
				errCodeString := string(errCode)
				if len(errCodeString) > 3 {
					c.JSON(http.StatusBadRequest, MSPNormalizeHTTPResponseBody{
						Code:    mspErrCodeTransform(errCodeString),
						Message: err.Error(),
					})
					return
				} else {
					errCodeInt, _ := strconv.Atoi(string(errCode))
					c.JSON(errCodeInt, MSPNormalizeHTTPResponseBody{
						Code:    mspErrCodeTransform(errCodeInt),
						Message: err.Error(),
					})
					return
				}
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
