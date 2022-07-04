package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/controlplane/dao/redis"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/service"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/dao/api"
	userDao "github.com/cloudslit/cloudslit/fullnode/app/v1/user/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	oauth2Help "github.com/cloudslit/cloudslit/fullnode/pkg/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func GetRedirectURL(c *gin.Context, company string) (redirectURL string, code int) {
	info, err := mysql.NewOauth2(c).GetOauth2ByCompany(company)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if info.ID == 0 {
		code = pconst.CODE_API_BAD_REQUEST
		return
	}
	config, err := service.Oauth2Config(info)
	if err != nil {
		code = pconst.CODE_API_BAD_REQUEST
		return
	}
	redirectURL, err = oauth2Help.GetOauth2RedirectURL(c, config)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func Oauth2Callback(c *gin.Context, session sessions.Session, company, oauth2Code string) {
	var user *mmysql.User
	// 查询对应的配置
	info, err := mysql.NewOauth2(c).GetOauth2ByCompany(company)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	if info.ID == 0 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	config, err := service.Oauth2Config(info)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	switch company {
	case "github":
		githubUser, err := api.GetGithubUser(c, config, oauth2Code)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
			return
		}
		user = &mmysql.User{Email: fmt.Sprintf("%s@github.com", *githubUser.Login), AvatarUrl: *githubUser.AvatarURL}
		if err = userDao.NewUser(c).FirstOrCreateUser(user); err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
			return
		}
	case "google":
		config.Endpoint = google.Endpoint
	case "facebook":
		config.Endpoint = facebook.Endpoint
	default:

	}
	userBytes, _ := json.Marshal(user)
	session.Set("user", userBytes)
	session.Save()
	// 判断是否有机器鉴权
	if machine := session.Get("machine"); machine != nil {
		// 删除machine这个session
		session.Delete("machine")
		session.Save()
		// 给当前请求授权
		cookie, _ := c.Cookie("zta")
		redis.NewMachine(c).PubMachineCookie(machine.(string), cookie)
		c.String(http.StatusOK, "Auth Success!")
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
	return
}
