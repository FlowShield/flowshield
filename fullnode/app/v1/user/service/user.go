package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mparam"

	"golang.org/x/oauth2/github"

	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
	"golang.org/x/oauth2"

	"github.com/google/uuid"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/controlplane/dao/redis"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/dao/api"
	userDao "github.com/cloudslit/cloudslit/fullnode/app/v1/user/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	oauth2Help "github.com/cloudslit/cloudslit/fullnode/pkg/oauth2"
	"github.com/gin-gonic/gin"
)

func GetRedirectURL(c *gin.Context) (redirectURL string, code int) {
	cfg := confer.GlobalConfig().Oauth2
	domain := confer.ConfigAppGetString("domain", "http://localhost:8080")
	if len(os.Getenv(domain)) > 0 {
		domain = os.Getenv(domain)
	}
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  domain + "/api/v1/user/oauth2/callback",
		Scopes:       []string{"user"},
	}
	redirectURL, err := oauth2Help.GetOauth2RedirectURL(c, config)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func Oauth2Callback(c *gin.Context, session sessions.Session, oauth2Code string) {
	var userInfo *mmysql.User
	cfg := confer.GlobalConfig().Oauth2
	domain := confer.ConfigAppGetString("domain", "http://localhost:8080")
	if len(os.Getenv(domain)) > 0 {
		domain = os.Getenv(domain)
	}
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  domain + "/api/v1/user/oauth2/callback",
		Scopes:       []string{"user"},
	}
	githubUser, err := api.GetGithubUser(c, config, oauth2Code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	user := &mmysql.User{
		Email:     fmt.Sprintf("%s@github.com", *githubUser.Login),
		AvatarUrl: *githubUser.AvatarURL,
		UUID:      uuid.NewString(),
	}
	if userInfo, err = userDao.NewUser(c).FirstOrCreateUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	// 判断是否Dao主
	userInfo.Master = userInfo.Wallet == confer.GlobalConfig().P2P.Account
	userBytes, _ := json.Marshal(userInfo)
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

func UserBindWallet(c *gin.Context, param *mparam.BindWallet) (code int) {
	// 查询是否已经绑定钱包
	user, err := userDao.NewUser(c).GetUser(util.User(c).UUID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if user.Wallet == param.Wallet {
		return pconst.CODE_COMMON_DATA_ALREADY_EXIST
	}
	user.Wallet = param.Wallet
	if err = userDao.NewUser(c).UpdateUser(user); err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	user.Master = user.Wallet == confer.GlobalConfig().P2P.Account
	// 绑定成功，刷新session
	session := sessions.Default(c)
	userBytes, _ := json.Marshal(user)
	session.Set("user", userBytes)
	session.Save()
	return
}
