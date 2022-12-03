package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/flowshield/flowshield/fullnode/app/v1/controlplane/dao/redis"
	"github.com/flowshield/flowshield/fullnode/app/v1/user/dao/api"
	userDao "github.com/flowshield/flowshield/fullnode/app/v1/user/dao/mysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/user/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	oauth2Help "github.com/flowshield/flowshield/fullnode/pkg/oauth2"
	"github.com/flowshield/flowshield/fullnode/pkg/util"
	"github.com/flowshield/flowshield/fullnode/pkg/web3/eth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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
	master, provider, err := eth.Instance().GetUserInfo(&bind.CallOpts{
		From: eth.CS.Auth.From,
	}, userInfo.UUID)
	if err != nil {
		logger.Errorf(c, "get wallet error: %v", err)
	} else {
		logger.Infof("get userinfo result: %v, master:%v, provider:%v", userInfo.UUID, master, provider)
		userInfo.Master = master
		userInfo.Provider = provider
	}
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

func UserRefresh(c *gin.Context) (code int) {
	// 查询信息
	user, err := userDao.NewUser(c).GetUser(util.User(c).UUID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	master, provider, err := eth.Instance().GetUserInfo(&bind.CallOpts{
		From: eth.CS.Auth.From,
	}, user.UUID)
	if err != nil {
		logger.Errorf(c, "get wallet error: %v", err)
	} else {
		logger.Infof("get userinfo result: %v, master:%v, provider:%v", user.UUID, master, provider)
		user.Master = master
		user.Provider = provider
	}
	//刷新session
	session := sessions.Default(c)
	userBytes, _ := json.Marshal(user)
	session.Set("user", userBytes)
	session.Save()
	return
}

func CheckAndBindUser(user *mmysql.User) (code int) {
	_, status, err := eth.Instance().GetWallet(&bind.CallOpts{
		From: eth.CS.Auth.From,
	}, user.UUID)
	if err != nil {
		logger.Errorf(nil, "contract get wallet error: %v", err)
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	//logger.Infof("CheckAndBindUse, uuid: %v,status: %v", user.UUID, status)
	switch status {
	//case 1:
	//	// 代表用户状态为预绑定，执行绑定
	//	tra, err := eth.Instance().VerifyWallet(eth.CS.Auth, user.UUID)
	//	if err != nil {
	//		logger.Errorf(nil, "contract verify wallet error: %v", err)
	//		return
	//	}
	//	rec, err := bind.WaitMined(context.Background(), eth.CS.Client, tra)
	//	if err != nil {
	//		logger.Errorf(nil, "contract verify wallet error: %v", err)
	//		return
	//	}
	//	if rec.Status == 0 {
	//		logger.Errorf(nil, "contract verify wallet err: %v", user.UUID)
	//		return
	//	}
	case 2:
	default:
		return
	}
	// 修改绑定状态并同步本地
	user.Status = mmysql.Bind
	err = userDao.NewUser(nil).UpdateUser(user)
	if err != nil {
		logger.Errorf(nil, "update user error: %v", err)
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}
