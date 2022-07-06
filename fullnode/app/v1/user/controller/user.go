package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mparam"

	"github.com/cloudslit/cloudslit/fullnode/app/base/controller"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/service"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	"github.com/cloudslit/cloudslit/fullnode/pkg/response"
	"github.com/cloudslit/cloudslit/fullnode/pkg/util"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	redirectURL, code := service.GetRedirectURL(c)
	if code != pconst.CODE_ERROR_OK {
		// TODO Redirect to BadRequest page
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("company %s not support", c.Param("company")))
		return
	}
	c.Redirect(http.StatusSeeOther, redirectURL)
}

func UserDetail(c *gin.Context) {
	response.UtilResponseReturnJson(c, pconst.CODE_ERROR_OK, util.User(c))
}

func Oauth2Callback(c *gin.Context) {
	session := sessions.Default(c)
	state := session.Get("state")
	if state != c.Query("state") {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("state error"))
		return
	}
	if len(c.Query("code")) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("code error"))
	}
	service.Oauth2Callback(c, session, c.Query("code"))
}

// @Summary UserBindWallet
// @Description 用户绑定钱包地址
// @Tags ZTA
// @Accept  json
// @Produce  json
// @Param Resource body mparam.AddResource true "用户绑定钱包地址"
// @Success 200 {object} controller.BindWallet
// @Router /user/bind [post]
func UserBindWallet(c *gin.Context) {
	param := &mparam.BindWallet{}
	b, code := controller.BindParams(c, &param)
	if !b {
		response.UtilResponseReturnJsonFailed(c, code)
		return
	}
	code = service.UserBindWallet(c, param)
	response.UtilResponseReturnJson(c, code, nil)
}
